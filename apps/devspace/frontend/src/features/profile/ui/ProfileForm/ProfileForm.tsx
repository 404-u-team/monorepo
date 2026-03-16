import type { JSX } from "react";
import { apiClient } from "@/shared/api/client";
import styles from "./ProfileForm.module.scss";
import { useEffect, useState } from "react";
import UserCard from "@/entities/user/ui/UserCard/UserCard";
import { Input, Button, Badge, Dropdown, Skeleton } from "@/shared/ui";
import { Camera, Save, X } from "lucide-react";
import { ProjectCard } from "@/entities/project";

//Вопрос 1: как запросить все скиллы, статусы и тд
//Вопрос 2: уведы где и когда и надо ли вообще
//Вопрос 3: в навыках в сваггере пока что нет цвета, делаем или нет

interface ProfileFormProps {
  id: string;
}

interface UserData {
  nickname: string;
  avatar_uri: string;
  main_role: string;
  bio: string;
  skills: string[];
  email: string;
}

const statusOptions = [
  { label: "В поиске", value: "searching" },
  { label: "Рассматриваю предложения", value: "considering" },
  { label: "Не ищу", value: "not_searching" },
];

const roleOptions = [
  { label: "Frontend Developer", value: "frontend" },
  { label: "Backend Developer", value: "backend" },
  { label: "Fullstack Developer", value: "fullstack" },
  { label: "UI/UX Designer", value: "designer" },
  { label: "Project Manager", value: "pm" },
];

const mockSkills = [
  { label: "React", value: "react", color: "3B82F6" },
  { label: "TypeScript", value: "ts", color: "34D399" },
  { label: "Node.js", value: "node", color: "F59E0B" },
  { label: "UI/UX", value: "uiux", color: "8B5CF6" },
];

export function ProfileForm({ id }: ProfileFormProps): JSX.Element {
  const [, setUserData] = useState<UserData | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [, setError] = useState<string | undefined>(undefined);

  const [nickname, setNickname] = useState("@maxik");
  const [email, setEmail] = useState("auth@gmail.com");
  const [status, setStatus] = useState("searching");
  const [role, setRole] = useState("frontend");
  const [bio, setBio] = useState("");
  const [searchSkill, setSearchSkill] = useState("");

  const project_id = "1";
  const avatar_uri =
    "https://img.freepik.com/premium-photo/vector-cat-with-character-wearing-jacket_575980-16303.jpg?semt=ais_hybrid";

  useEffect(() => {
    const fetchUserData = async (): Promise<void> => {
      try {
        setLoading(true);
        const response = await apiClient.get<UserData>(`/users/${id}`);
        setUserData(response.data);
        setError(undefined);
      } catch (error_) {
        setError("Ошибка загрузки данных пользователя");
        console.error("Error fetching user data:", error_);
      } finally {
        setLoading(false);
      }
    };
    void fetchUserData();
  }, [id]);

  if (loading) {
    return (
      <div className={styles.profileGrid}>
        <div className={styles.column1}>
          <Skeleton height={200} borderRadius={16} />
          <Skeleton height={120} borderRadius={16} />
        </div>
        <div className={styles.column2}>
          <Skeleton height={40} width="60%" />
          <Skeleton height={30} width="40%" />
          <Skeleton height={200} borderRadius={16} />
        </div>
      </div>
    );
  }

  return (
    <div className={styles.profileGrid}>
      <div className={styles.column1}>
        <h1>Мой Профиль</h1>
        <div className={styles.miniProfile}>
          <UserCard id={id} />
        </div>
        <div className={styles.pov}>
          <Button variant="primary" className={styles.profileButton}>
            Профиль
          </Button>
          <Button variant="outline" className={styles.notificationButton}>
            Уведомления
          </Button>
        </div>
      </div>

      <div className={styles.column2}>
        <h2>Настройки профиля</h2>

        <div className={styles.avatarChange}>
          <img src={avatar_uri} alt="Avatar" className={styles.avatar} />
          <h3>Редактировать изображение профиля</h3>
          <Button variant="outline" className={styles.changeAvaButton}>
            <Camera size={16} />
            Сменить фото
          </Button>
        </div>

        <div className={styles.mySkills}>
          <h2>Мои навыки</h2>
          <Input
            placeholder="Поиск навыков..."
            value={searchSkill}
            onChange={(event_) => {
              setSearchSkill(event_.target.value);
            }}
            className={styles.searchBox}
          />

          <div className={styles.skillShowcase}>
            {mockSkills.map((skill) => (
              <Badge key={skill.value} color={skill.color}>
                {skill.label}
              </Badge>
            ))}
          </div>

          <div className={styles.saveButtons}>
            <Button variant="primary">
              <Save size={25} />
              Сохранить изменения
            </Button>
            <Button variant="outline">
              <X size={16} />
              Отмена
            </Button>
          </div>
        </div>
      </div>

      <div className={styles.column3}>
        <div className={styles.settingsGrid}>
          <div className={styles.settingBox}>
            <h3>Никнейм</h3>
            <Input
              placeholder="Например: @maxik"
              value={nickname}
              onChange={(event_) => {
                setNickname(event_.target.value);
              }}
              className={styles.settingInput}
            />
          </div>

          <div className={styles.settingBox}>
            <h3>Статус</h3>
            <Dropdown
              options={statusOptions}
              value={status}
              onChange={setStatus}
              placeholder="Выберите статус"
            />
          </div>

          <div className={styles.settingBox}>
            <h3>E-mail</h3>
            <Input
              type="email"
              placeholder="Например: auth@gmail.com"
              value={email}
              onChange={(event_) => {
                setEmail(event_.target.value);
              }}
              className={styles.settingInput}
            />
          </div>

          <div className={styles.settingBox}>
            <h3>Основная роль</h3>
            <Dropdown
              options={roleOptions}
              value={role}
              onChange={setRole}
              placeholder="Выберите роль"
            />
          </div>
        </div>
        <div className={styles.bio}>
          <h3>Описание</h3>
          <textarea
            className={styles.bioInput}
            value={bio}
            onChange={(event_) => {
              setBio(event_.target.value);
            }}
            placeholder="Расскажите о себе..."
            rows={4}
          />
        </div>
        {/* <div className={styles.miniNav}>
          <h3>Быстрая навигация</h3>
          <Button variant="clear">Проекты</Button>
          <Button variant="clear">Команды</Button>
          <Button variant="clear">Достижения</Button>
        </div> */}

        <div className={styles.feed}>
          <h3>Лента активности</h3>
          {/* <p className={styles.emptyFeed}>Здесь пока ничего нет</p> */}
          <ProjectCard projectId={project_id} />
        </div>
      </div>
    </div>
  );
}
