import { useNavigate } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useState, type JSX } from "react";

import { createIdea } from "@/entities/idea";
import { Button, Input, Badge, MdEditor } from "@/shared/ui";

import styles from "./CreateIdeaForm.module.scss";

export const CreateIdeaForm = observer(function CreateIdeaForm(): JSX.Element {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [category, setCategory] = useState("");
  const [content, setContent] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | undefined>(undefined);

  const handleSubmit = async (event: React.SyntheticEvent): Promise<void> => {
    event.preventDefault();
    if (!title) return;

    setIsSubmitting(true);
    setError(undefined);
    try {
      const newIdea = await createIdea({ title, description, content, category });
      void navigate({ to: `/idea/${newIdea.id}` });
    } catch {
      setError("Произошла ошибка при создании идеи. Попробуйте снова.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.formSection}>
        <h1 className={styles.title}>Добавить новую идею</h1>
        <form
          className={styles.form}
          onSubmit={(event) => {
            void handleSubmit(event);
          }}
        >
          <div className={styles.field}>
            <label className={styles.label} htmlFor="idea-title">
              Название идеи
            </label>
            <Input
              id="idea-title"
              placeholder="Введите название..."
              value={title}
              onChange={(event_) => {
                setTitle(event_.target.value);
              }}
              required
              disabled={isSubmitting}
            />
          </div>
          <div className={styles.field}>
            <label className={styles.label} htmlFor="idea-category">
              Категория
            </label>
            <Input
              id="idea-category"
              placeholder="Например: Education, Technology..."
              value={category}
              onChange={(event_) => {
                setCategory(event_.target.value);
              }}
              disabled={isSubmitting}
            />
          </div>
          <div className={styles.field}>
            <label className={styles.label} htmlFor="idea-description">
              Краткое описание
            </label>
            <textarea
              id="idea-description"
              className={styles.textarea}
              placeholder="Краткое описание вашей идеи..."
              value={description}
              onChange={(event_) => {
                setDescription(event_.target.value);
              }}
              disabled={isSubmitting}
            />
          </div>
          <div className={styles.field}>
            <label className={styles.label} htmlFor="idea-content">
              Содержимое
            </label>
            <MdEditor
              id="idea-content"
              value={content}
              onChange={setContent}
              placeholder="Напишите подробное содержимое идеи..."
              height={360}
              disabled={isSubmitting}
            />
          </div>

          {error !== undefined && <p className={styles.error}>{error}</p>}

          <Button type="submit" fullWidth disabled={isSubmitting || !title}>
            <Plus size={18} />
            Опубликовать идею
          </Button>
        </form>
      </div>

      <div className={styles.previewSection}>
        <h2 className={styles.previewTitle}>Предпросмотр карточки</h2>
        <div className={styles.previewCard}>
          <article className={styles.mockCard}>
            <div className={styles.imagePlaceholder}>
              <div className={styles.gradient} />
              {category !== "" && <Badge className={styles.badge}>{category}</Badge>}
            </div>
            <div className={styles.cardBody}>
              <h3 className={styles.cardTitle}>{title || "Заголовок идеи"}</h3>
              <p className={styles.cardDescription}>
                {description || "Здесь будет описание вашей идеи..."}
              </p>
            </div>
            <div className={styles.cardFooter}>
              <div className={styles.mockUser}>
                <div className={styles.mockAvatar} />
                <span className={styles.mockNickname}>Вы</span>
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>
  );
});
