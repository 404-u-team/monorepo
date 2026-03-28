import { useState, type JSX } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { observer } from 'mobx-react-lite';
import { Plus } from 'lucide-react';
import { Button, Input } from '@/shared/ui';
import { createProject } from '@/entities/project';
import styles from './CreateProjectForm.module.scss';

interface CreateProjectFormProps {
    initialTitle?: string | undefined;
    initialDescription?: string | undefined;
    initialIdeaId?: string | undefined;
}

export const CreateProjectForm = observer(function CreateProjectForm({
    initialTitle = '',
    initialDescription = '',
    initialIdeaId,
}: CreateProjectFormProps): JSX.Element {
    const navigate = useNavigate();
    const [title, setTitle] = useState(initialTitle);
    const [description, setDescription] = useState(initialDescription);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | undefined>(undefined);

    const handleSubmit = async (event: React.SyntheticEvent): Promise<void> => {
        event.preventDefault();
        if (!title) return;

        setIsSubmitting(true);
        setError(undefined);
        try {
            const newProject = await createProject({ title, description, idea_id: initialIdeaId });
            void navigate({ to: `/project/$projectId`, params: { projectId: newProject.id } });
        } catch {
            setError('Произошла ошибка при создании проекта. Попробуйте снова.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.formSection}>
                <h1 className={styles.title}>Создать новый проект</h1>
                <form
                    className={styles.form}
                    onSubmit={(event) => { void handleSubmit(event); }}
                >
                    <div className={styles.field}>
                        <label className={styles.label} htmlFor="project-title">Название проекта</label>
                        <Input
                            id="project-title"
                            placeholder="Введите название..."
                            value={title}
                            onChange={(event_) => { setTitle(event_.target.value); }}
                            required
                            disabled={isSubmitting}
                        />
                    </div>
                    <div className={styles.field}>
                        <label className={styles.label} htmlFor="project-description">Описание</label>
                        <textarea
                            id="project-description"
                            className={styles.textarea}
                            placeholder="Опишите цели и задачи проекта..."
                            value={description}
                            onChange={(event_) => { setDescription(event_.target.value); }}
                            disabled={isSubmitting}
                        />
                    </div>

                    {error !== undefined && <p className={styles.error}>{error}</p>}

                    <Button type="submit" fullWidth disabled={isSubmitting || !title}>
                        <Plus size={18} />
                        Создать проект
                    </Button>
                </form>
            </div>

            <div className={styles.previewSection}>
                <h2 className={styles.previewTitle}>Предпросмотр карточки</h2>
                <div className={styles.previewCard}>
                    <article className={styles.mockCard}>
                        <div className={styles.imagePlaceholder}>
                            <div className={styles.gradient} />
                        </div>
                        <div className={styles.cardBody}>
                            <h3 className={styles.cardTitle}>{title || 'Название проекта'}</h3>
                            <p className={styles.cardDescription}>
                                {description || 'Здесь будет описание вашего проекта...'}
                            </p>
                        </div>
                        <div className={styles.cardFooter}>
                            <div className={styles.mockUser}>
                                <div className={styles.mockAvatar} />
                                <span className={styles.mockNickname}>Вы — лидер</span>
                            </div>
                        </div>
                    </article>
                </div>
            </div>
        </div>
    );
});
