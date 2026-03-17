import { useState, type JSX } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { observer } from 'mobx-react-lite';
import { Plus } from 'lucide-react';
import { Button, Input, Badge } from '@/shared/ui';
import { createIdea } from '@/entities/idea';
import styles from './CreateIdeaForm.module.scss';

export const CreateIdeaForm = observer(function CreateIdeaForm(): JSX.Element {
    const navigate = useNavigate();
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | undefined>(undefined);

    const handleSubmit = async (event: React.FormEvent): Promise<void> => {
        event.preventDefault();
        if (!title || !description) return;

        setIsSubmitting(true);
        setError(undefined);
        try {
            const newIdea = await createIdea({ title, description });
            void navigate({ to: `/idea/${newIdea.id}` });
        } catch {
            setError('Произошла ошибка при создании идеи. Попробуйте снова.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.formSection}>
                <h1 className={styles.title}>Добавить новую идею</h1>
                <form className={styles.form} onSubmit={handleSubmit}>
                    <div className={styles.textareaWrapper}>
                        <label className={styles.label}>Название идеи</label>
                        <Input
                            placeholder="Введите название..."
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            required
                            disabled={isSubmitting}
                        />
                    </div>
                    <div className={styles.textareaWrapper}>
                        <label className={styles.label}>Описание</label>
                        <textarea
                            className={styles.textarea}
                            placeholder="Опишите вашу идею подробно..."
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            required
                            disabled={isSubmitting}
                        />
                    </div>

                    {error !== undefined && <p className={styles.error}>{error}</p>}

                    <Button type="submit" fullWidth disabled={isSubmitting || !title || !description}>
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
                            <Badge className={styles.badge}>Preview</Badge>
                        </div>
                        <div className={styles.cardBody}>
                            <h3 className={styles.cardTitle}>{title || 'Заголовок идеи'}</h3>
                            <p className={styles.cardDescription}>
                                {description || 'Здесь будет описание вашей идеи...'}
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
