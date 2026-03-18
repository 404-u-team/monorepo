import type { JSX } from "react";
import { apiClient } from "@/shared/api/client";
import styles from "./ProfileForm.module.scss";
import { useEffect, useState } from "react";
import UserCard from "@/entities/user/ui/UserCard/UserCard";
import { Input, Button, Badge, Dropdown, Skeleton } from "@/shared/ui";
import { Camera, Save, X } from "lucide-react";
import { ProjectCard } from "@/entities/project";
import { statusOptions, roleOptions } from "@/shared/enums/ProfileEnums";

interface ProfileFormProps {
  id: string;
}

interface UserData {
  id: string;
  nickname: string;
  avatar_uri: string;
  main_role: string;
  bio: string;
  skills: string[];
  email: string;
  status?: string;
}

interface UpdateUserData {
  nickname?: string;
  bio?: string;
}

interface Skill {
  label: string;
  value: string;
  color: string;
}

//Временные моки, пока API не готово
const mockSkills = [
  { label: "React", value: "react", color: "3B82F6" },
  { label: "TypeScript", value: "ts", color: "34D399" },
  { label: "Node.js", value: "node", color: "F59E0B" },
  { label: "UI/UX", value: "uiux", color: "8B5CF6" },
];

export function ProfileForm({ id }: ProfileFormProps): JSX.Element {
  const [userData, setUserData] = useState<UserData>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | undefined>(undefined);
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [nickname, setNickname] = useState("");
  const [email, setEmail] = useState("");
  const [status, setStatus] = useState("searching");
  const [role, setRole] = useState("frontend");
  const [bio, setBio] = useState("");
  const [searchSkill, setSearchSkill] = useState("");
  const [initialNickname, setInitialNickname] = useState("");
  const [initialBio, setInitialBio] = useState("");

  const avatar_uri =
    userData?.avatar_uri ??
    "https://img.freepik.com/premium-photo/vector-cat-with-character-wearing-jacket_575980-16303.jpg?semt=ais_hybrid";

  const project_id = "1";

  useEffect(() => {
    const fetchUserData = async (): Promise<void> => {
      try {
        setLoading(true);
        setError(undefined);
        const endpoint = `/users/me`;
        const response = await apiClient.get<UserData>(endpoint);
        const data = response.data;
        setUserData(data);
        setNickname(data.nickname || "");
        setEmail(data.email || "");
        setBio(data.bio || "");
        setRole(data.main_role || "frontend");
        //setStatus(data.status || "searching");

        setInitialNickname(data.nickname || "");
        setInitialBio(data.bio || "");
      } catch (error_) {
        setError("Ошибка загрузки данных пользователя");
        console.error("Error fetching user data:", error_);
      } finally {
        setLoading(false);
      }
    };

    void fetchUserData();
  }, [id]);

  const handleSaveProfile = async (): Promise<void> => {
    const hasNicknameChanged = nickname !== initialNickname;
    const hasBioChanged = bio !== initialBio;

    if (!hasNicknameChanged && !hasBioChanged) {
      return;
    }

    try {
      setIsSaving(true);
      setError(undefined);
      setSaveSuccess(false);

      const updateData: UpdateUserData = {};

      if (hasNicknameChanged) {
        updateData.nickname = nickname;
      }
      if (hasBioChanged) {
        updateData.bio = bio;
      }
      await apiClient.put("/users/me", updateData);
      setInitialNickname(nickname);
      setInitialBio(bio);
      setSaveSuccess(true);
      setTimeout(() => {
        setSaveSuccess(false);
      }, 3000);
    } catch (error_) {
      setError("Ошибка при сохранении изменений");
      console.error("Error saving user data:", error_);
    } finally {
      setIsSaving(false);
    }
  };
  const handleCancelChanges = (): void => {
    setNickname(initialNickname);
    setBio(initialBio);
  };

  const hasChanges = nickname !== initialNickname || bio !== initialBio;

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

  if (error !== undefined && !userData) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.error}>{error}</p>
        <Button
          variant="primary"
          onClick={() => {
            window.location.reload();
          }}
        >
          Попробовать снова
        </Button>
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

        {saveSuccess && (
          <div className={styles.successMessage}>
            Изменения успешно сохранены!
          </div>
        )}

        {error !== undefined && (
          <div className={styles.errorMessage}>{error}</div>
        )}

        <div className={styles.avatarChange}>
          <div className={styles.avatar}>
            <img src={avatar_uri} alt="Avatar" />
          </div>
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
            {mockSkills.map((skill: Skill) => (
              <Badge key={skill.value} color={skill.color}>
                {skill.label}
                <X size={16} />
              </Badge>
            ))}
          </div>

          <div className={styles.saveButtons}>
            <Button
              variant="primary"
              onClick={() => void handleSaveProfile()}
              disabled={isSaving || !hasChanges}
            >
              <Save size={25} />
              {isSaving ? "Сохранение..." : "Сохранить изменения"}
            </Button>
            <Button
              variant="outline"
              onClick={handleCancelChanges}
              disabled={isSaving || !hasChanges}
            >
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
              disabled
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

        <div className={styles.feed}>
          <h3>Лента активности</h3>
          <ProjectCard projectId={project_id} />
        </div>
      </div>
    </div>
  );
}
