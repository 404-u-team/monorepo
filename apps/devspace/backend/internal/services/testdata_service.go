package services

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tdTotalSkills   = 800
	tdRootSkills    = 80
	tdTotalUsers    = 500
	tdTotalIdeas    = 1000
	tdTotalProjects = 1000
	// grand total for progress calculation (slots are bonus, not counted)
	tdTotalItems = tdTotalSkills + tdTotalUsers + tdTotalIdeas + tdTotalProjects
)

// GenerationStatus holds the current state of the async generation job.
type GenerationStatus struct {
	Running  bool   `json:"running"`
	Stage    string `json:"stage"`
	Progress int    `json:"progress"` // 0–100
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

// TestDataService manages async test-data generation.
type TestDataService interface {
	Start() error
	Cancel()
	Status() GenerationStatus
}

type testDataService struct {
	db     *gorm.DB
	cfg    *config.Config
	mu     sync.RWMutex
	status GenerationStatus
	cancel context.CancelFunc
}

func NewTestDataService(db *gorm.DB, cfg *config.Config) TestDataService {
	return &testDataService{db: db, cfg: cfg}
}

func (s *testDataService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.status.Running {
		return fmt.Errorf("генерация данных уже запущена")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.status = GenerationStatus{Running: true, Stage: "Инициализация", Progress: 0}
	go s.generate(ctx)
	return nil
}

func (s *testDataService) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *testDataService) Status() GenerationStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *testDataService) setStage(stage string) {
	s.mu.Lock()
	s.status.Stage = stage
	s.mu.Unlock()
}

func (s *testDataService) setProgress(completed int) {
	s.mu.Lock()
	s.status.Progress = completed * 100 / tdTotalItems
	s.mu.Unlock()
}

func (s *testDataService) finish(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Running = false
	if err != nil {
		s.status.Stage = "Ошибка"
		s.status.Error = err.Error()
	} else {
		s.status.Stage = "Завершено"
		s.status.Progress = 100
		s.status.Done = true
	}
}

// ─── main goroutine ──────────────────────────────────────────────────────────

func (s *testDataService) generate(ctx context.Context) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	completed := 0

	// Pre-hash a single password once and reuse it for all test users.
	hashedPassword, err := auth.HashPassword("TestPass123!", s.cfg)
	if err != nil {
		s.finish(fmt.Errorf("не удалось хэшировать пароль: %w", err))
		return
	}

	// ── 1. Skills ────────────────────────────────────────────────────────────
	s.setStage("Генерация навыков")

	rootNames := tdRootSkillNames()
	rootIDs := make([]uuid.UUID, tdRootSkills)
	rootModels := make([]models.SkillCategory, tdRootSkills)
	for i := 0; i < tdRootSkills; i++ {
		id := uuid.New()
		rootIDs[i] = id
		rootModels[i] = models.SkillCategory{ID: id, Name: rootNames[i]}
	}
	if err := s.db.WithContext(ctx).CreateInBatches(rootModels, 100).Error; err != nil {
		s.finish(err)
		return
	}
	completed += tdRootSkills
	s.setProgress(completed)

	childCount := tdTotalSkills - tdRootSkills // 720
	childPerRoot := childCount / tdRootSkills  // 9
	childModels := make([]models.SkillCategory, 0, childCount)
	for i, parentID := range rootIDs {
		if ctx.Err() != nil {
			s.finish(fmt.Errorf("отменено"))
			return
		}
		pid := parentID
		for j := 0; j < childPerRoot; j++ {
			childModels = append(childModels, models.SkillCategory{
				ID:       uuid.New(),
				ParentID: &pid,
				Name:     fmt.Sprintf("%s / Навык %d", rootNames[i], j+1),
			})
		}
	}
	if err := s.db.WithContext(ctx).CreateInBatches(childModels, 200).Error; err != nil {
		s.finish(err)
		return
	}
	completed += childCount
	s.setProgress(completed)

	allSkillIDs := make([]uuid.UUID, 0, tdTotalSkills)
	allSkillIDs = append(allSkillIDs, rootIDs...)
	for _, sc := range childModels {
		allSkillIDs = append(allSkillIDs, sc.ID)
	}

	childrenByRoot := make(map[uuid.UUID][]uuid.UUID, len(rootIDs))
	for _, rootID := range rootIDs {
		childrenByRoot[rootID] = []uuid.UUID{}
	}
	for _, sc := range childModels {
		if sc.ParentID != nil {
			childrenByRoot[*sc.ParentID] = append(childrenByRoot[*sc.ParentID], sc.ID)
		}
	}

	// ── 2. Users ─────────────────────────────────────────────────────────────
	s.setStage("Генерация пользователей")

	bios := []string{
		"Full-stack разработчик с опытом 3+ лет",
		"Люблю создавать качественные продукты",
		"Специализируюсь на backend-разработке",
		"Увлекаюсь ML и анализом данных",
		"DevOps-инженер, оркестрация контейнеров",
		"Frontend-разработчик, React и Vue",
		"Студент, изучаю Go и микросервисы",
		"Серийный предприниматель в IT",
		"Архитектор программных систем",
		"Мобильная разработка iOS/Android",
	}

	userIDs := make([]uuid.UUID, tdTotalUsers)
	userModels := make([]models.User, tdTotalUsers)
	for i := 0; i < tdTotalUsers; i++ {
		id := uuid.New()
		userIDs[i] = id
		mainRoleID := rootIDs[rng.Intn(len(rootIDs))]
		userModels[i] = models.User{
			ID:           id,
			Email:        fmt.Sprintf("testuser_%04d@devspace.test", i+1),
			Nickname:     fmt.Sprintf("testuser_%04d", i+1),
			PasswordHash: hashedPassword,
			Bio:          bios[i%len(bios)],
			MainRole:     &mainRoleID,
		}
	}
	if err := s.db.WithContext(ctx).CreateInBatches(userModels, 100).Error; err != nil {
		s.finish(err)
		return
	}
	completed += tdTotalUsers
	s.setProgress(completed)

	// Assign 1–3 skills to ~50 % of users.
	userSkillModels := make([]models.UserSkill, 0, tdTotalUsers*2)
	for i := 0; i < tdTotalUsers/2; i++ {
		if ctx.Err() != nil {
			s.finish(fmt.Errorf("отменено"))
			return
		}
		userID := userIDs[i]
		used := make(map[uuid.UUID]bool)
		for k := 0; k < rng.Intn(3)+1; k++ {
			skillID := allSkillIDs[rng.Intn(len(allSkillIDs))]
			if !used[skillID] {
				used[skillID] = true
				userSkillModels = append(userSkillModels, models.UserSkill{UserID: userID, SkillID: skillID})
			}
		}
	}
	if len(userSkillModels) > 0 {
		if err := s.db.WithContext(ctx).CreateInBatches(userSkillModels, 200).Error; err != nil {
			s.finish(err)
			return
		}
	}

	// ── 3. Ideas ─────────────────────────────────────────────────────────────
	s.setStage("Генерация идей")

	ideaCategories := []string{
		"Education", "Healthcare", "Finance", "Entertainment", "E-commerce",
		"Social", "Gaming", "Productivity", "Travel", "Agriculture",
	}
	ideaDescriptions := []string{
		"Платформа для автоматизации рабочих процессов с применением AI.",
		"Мобильное приложение для удобного управления повседневными задачами.",
		"Сервис для поиска специалистов в различных профессиональных областях.",
		"Образовательная платформа с обучением через практические проекты.",
		"Система мониторинга и аналитики для малого и среднего бизнеса.",
		"Приложение для соединения волонтёров с некоммерческими организациями.",
		"Маркетплейс для продажи цифровых товаров и авторского контента.",
		"Сервис планирования путешествий с персонализированными рекомендациями.",
		"Инструмент для эффективной совместной работы распределённых команд.",
		"Платформа обмена опытом между студентами и практикующими специалистами.",
	}

	ideaIDs := make([]uuid.UUID, tdTotalIdeas)
	ideaModels := make([]models.Idea, tdTotalIdeas)
	for i := 0; i < tdTotalIdeas; i++ {
		if ctx.Err() != nil {
			s.finish(fmt.Errorf("отменено"))
			return
		}
		id := uuid.New()
		ideaIDs[i] = id
		short := id.String()[:8]
		ideaModels[i] = models.Idea{
			ID:          id,
			AuthorID:    userIDs[rng.Intn(len(userIDs))],
			Title:       fmt.Sprintf("[%s] %s #%d", short, ideaCategories[i%len(ideaCategories)], i+1),
			Description: ideaDescriptions[i%len(ideaDescriptions)],
			Category:    ideaCategories[i%len(ideaCategories)],
		}
	}
	if err := s.db.WithContext(ctx).CreateInBatches(ideaModels, 200).Error; err != nil {
		s.finish(err)
		return
	}
	completed += tdTotalIdeas
	s.setProgress(completed)

	// ── 4. Projects ───────────────────────────────────────────────────────────
	s.setStage("Генерация проектов")

	projectTitles := []string{
		"AI Language Learning", "Smart City Platform", "HealthTracker Pro",
		"EduConnect", "FinanceBot", "TravelMate", "AgriSmart", "GameWorld",
		"SocialHub", "ProductivitySuite", "CodeReview AI", "DevMetrics",
		"OpenMarket", "CloudDeploy", "SecureVault",
	}
	projectDescriptions := []string{
		"Разработка инновационного продукта для широкой аудитории пользователей.",
		"Платформа нового поколения для автоматизации бизнес-процессов.",
		"Мобильное приложение с AI-функциями для повседневного использования.",
		"B2B SaaS-решение для оптимизации рабочих процессов внутри команд.",
		"Образовательная платформа с геймификацией и персонализацией обучения.",
	}
	statuses := []string{"open", "closed"}

	projectIDs := make([]uuid.UUID, tdTotalProjects)
	projectModels := make([]models.Project, tdTotalProjects)
	for i := 0; i < tdTotalProjects; i++ {
		if ctx.Err() != nil {
			s.finish(fmt.Errorf("отменено"))
			return
		}
		id := uuid.New()
		projectIDs[i] = id
		short := id.String()[:8]
		leaderID := userIDs[rng.Intn(len(userIDs))]

		var ideaID *uuid.UUID
		// first 25 % of projects are linked to an idea
		if i < tdTotalProjects/4 {
			idea := ideaIDs[rng.Intn(len(ideaIDs))]
			ideaID = &idea
		}

		desc := projectDescriptions[i%len(projectDescriptions)]
		projectModels[i] = models.Project{
			ID:          id,
			LeaderID:    leaderID,
			IdeaID:      ideaID,
			Title:       fmt.Sprintf("%s [%s]", projectTitles[i%len(projectTitles)], short),
			Description: &desc,
			Status:      statuses[rng.Intn(len(statuses))],
		}
	}
	if err := s.db.WithContext(ctx).CreateInBatches(projectModels, 200).Error; err != nil {
		s.finish(err)
		return
	}
	completed += tdTotalProjects
	s.setProgress(completed)

	// ── 5. Slots ─────────────────────────────────────────────────────────────
	s.setStage("Генерация слотов")

	slotTitles := []string{
		"Backend Developer", "Frontend Developer", "Mobile Developer",
		"DevOps Engineer", "Data Scientist", "ML Engineer",
		"QA Engineer", "UI/UX Designer", "Product Manager", "Tech Lead",
	}
	// track used (projectID → set of userIDs) to satisfy the partial unique index
	projectUsedUsers := make(map[uuid.UUID]map[uuid.UUID]bool, tdTotalProjects)
	for _, pid := range projectIDs {
		projectUsedUsers[pid] = make(map[uuid.UUID]bool)
	}

	slotModels := make([]models.ProjectSlot, 0, tdTotalProjects*2)
	for i, projectID := range projectIDs {
		if ctx.Err() != nil {
			s.finish(fmt.Errorf("отменено"))
			return
		}
		// Distribution:
		//  0–299  (30 %) → no slots
		// 300–599 (30 %) → 2–3 open slots (no user assigned)
		// 600–799 (20 %) → 1–2 open + 1–2 filled slots
		// 800–999 (20 %) → 2–3 filled slots (user assigned, status closed)
		switch {
		case i < 300:
			// no slots

		case i < 600:
			for j := 0; j < rng.Intn(2)+2; j++ {
				primarySkillID, secondarySkillID := pickValidPrimarySecondarySkillPair(rng, rootIDs, childrenByRoot)
				slotModels = append(slotModels, models.ProjectSlot{
					ID:                uuid.New(),
					ProjectID:         projectID,
					PrimarySkillsID:   models.UUIDArray{primarySkillID},
					SecondarySkillsID: models.UUIDArray{secondarySkillID},
					Title:             slotTitles[rng.Intn(len(slotTitles))],
					Status:            "open",
				})
			}

		case i < 800:
			for j := 0; j < rng.Intn(2)+1; j++ {
				primarySkillID, secondarySkillID := pickValidPrimarySecondarySkillPair(rng, rootIDs, childrenByRoot)
				slotModels = append(slotModels, models.ProjectSlot{
					ID:                uuid.New(),
					ProjectID:         projectID,
					PrimarySkillsID:   models.UUIDArray{primarySkillID},
					SecondarySkillsID: models.UUIDArray{secondarySkillID},
					Title:             slotTitles[rng.Intn(len(slotTitles))],
					Status:            "open",
				})
			}
			for j := 0; j < rng.Intn(2)+1; j++ {
				userID := pickUniqueUser(rng, userIDs, projectUsedUsers[projectID])
				if userID == uuid.Nil {
					break
				}
				primarySkillID, secondarySkillID := pickValidPrimarySecondarySkillPair(rng, rootIDs, childrenByRoot)
				projectUsedUsers[projectID][userID] = true
				slotModels = append(slotModels, models.ProjectSlot{
					ID:                uuid.New(),
					ProjectID:         projectID,
					PrimarySkillsID:   models.UUIDArray{primarySkillID},
					SecondarySkillsID: models.UUIDArray{secondarySkillID},
					Title:             slotTitles[rng.Intn(len(slotTitles))],
					Status:            "closed",
					UserID:            &userID,
				})
			}

		default:
			for j := 0; j < rng.Intn(2)+2; j++ {
				userID := pickUniqueUser(rng, userIDs, projectUsedUsers[projectID])
				if userID == uuid.Nil {
					break
				}
				primarySkillID, secondarySkillID := pickValidPrimarySecondarySkillPair(rng, rootIDs, childrenByRoot)
				projectUsedUsers[projectID][userID] = true
				slotModels = append(slotModels, models.ProjectSlot{
					ID:                uuid.New(),
					ProjectID:         projectID,
					PrimarySkillsID:   models.UUIDArray{primarySkillID},
					SecondarySkillsID: models.UUIDArray{secondarySkillID},
					Title:             slotTitles[rng.Intn(len(slotTitles))],
					Status:            "closed",
					UserID:            &userID,
				})
			}
		}
	}

	if len(slotModels) > 0 {
		if err := s.db.WithContext(ctx).CreateInBatches(slotModels, 200).Error; err != nil {
			s.finish(err)
			return
		}
	}

	s.finish(nil)
}

// pickUniqueUser returns a userID not yet assigned to the given project.
// Returns uuid.Nil if no free user is found after reasonable attempts.
func pickUniqueUser(rng *rand.Rand, userIDs []uuid.UUID, used map[uuid.UUID]bool) uuid.UUID {
	if len(used) >= len(userIDs) {
		return uuid.Nil
	}
	for range 20 {
		id := userIDs[rng.Intn(len(userIDs))]
		if !used[id] {
			return id
		}
	}
	// linear fallback
	for _, id := range userIDs {
		if !used[id] {
			return id
		}
	}
	return uuid.Nil
}

func pickValidPrimarySecondarySkillPair(rng *rand.Rand, rootIDs []uuid.UUID, childrenByRoot map[uuid.UUID][]uuid.UUID) (uuid.UUID, uuid.UUID) {
	for range 20 {
		primary := rootIDs[rng.Intn(len(rootIDs))]
		children := childrenByRoot[primary]
		if len(children) == 0 {
			continue
		}
		secondary := children[rng.Intn(len(children))]
		return primary, secondary
	}

	// deterministic fallback; with generated data this branch should be unreachable.
	primary := rootIDs[0]
	children := childrenByRoot[primary]
	if len(children) == 0 {
		return primary, primary
	}
	return primary, children[0]
}

// tdRootSkillNames returns 80 unique root skill category names.
func tdRootSkillNames() []string {
	bases := []string{
		"Backend Development", "Frontend Development", "Mobile Development",
		"DevOps", "Machine Learning", "Data Science", "Cloud Computing",
		"Database Administration", "Cybersecurity", "Testing & QA",
		"UI/UX Design", "Product Management", "Game Development",
		"Blockchain", "Embedded Systems", "AR/VR Development",
		"Technical Writing", "Scrum & Agile", "System Architecture", "Open Source",
	}
	names := make([]string, tdRootSkills)
	for i := range tdRootSkills {
		base := bases[i%len(bases)]
		iteration := i/len(bases) + 1
		if iteration == 1 {
			names[i] = base
		} else {
			names[i] = fmt.Sprintf("%s (группа %d)", base, iteration)
		}
	}
	return names
}
