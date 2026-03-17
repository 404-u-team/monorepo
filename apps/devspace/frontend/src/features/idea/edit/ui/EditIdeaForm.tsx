import { useState, type JSX } from 'react';
import { observer } from 'mobx-react-lite';
import { Button, Input, MdEditor } from '@/shared/ui';
import { updateIdea, type IIdea } from '@/entities/idea';
import styles from './EditIdeaForm.module.scss';

interface EditIdeaFormProps {
    idea: IIdea;
    onSuccess: (updated: IIdea) => void;
    onCancel: () => void;
}

export const EditIdeaForm = observer(function EditIdeaForm({ idea, onSuccess, onCancel }: EditIdeaFormProps): JSX.Element {
    const [title, setTitle] = useState(idea.title);
    const [description, setDescription] = useState(idea.description);
    const [category, setCategory] = useState(idea.category ?? '');
    const [content, setContent] = useState(idea.content ?? '');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | undefined>(undefined);

    const handleSubmit = async (event: React.SyntheticEvent): Promise<void> => {
        event.preventDefault();
        if (!title) return;

        setIsSubmitting(true);
        setError(undefined);
        try {
            const updated = await updateIdea(idea.id, { title, description, content, category });
            onSuccess(updated);
        } catch {
            setError('Произошла ошибка при сохранении. Попробуйте снова.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <form 
            className={styles.form} 
            onSubmit={(event) => { void handleSubmit(event); }}
        >
            <div className={styles.field}>
                <label className={styles.label}>Название</label>
                <Input
                    value={title}
                    onChange={(event_) => { setTitle(event_.target.value); }}
                    placeholder="Введите название..."
                    required
                    disabled={isSubmitting}
                />
            </div>
            <div className={styles.field}>
                <label className={styles.label}>Категория</label>
                <Input
                    value={category}
                    onChange={(event_) => { setCategory(event_.target.value); }}
                    placeholder="Например: Education, Technology..."
                    disabled={isSubmitting}
                />
            </div>
            <div className={styles.field}>
                <label className={styles.label}>Краткое описание</label>
                <textarea
                    className={styles.textarea}
                    value={description}
                    onChange={(event_) => { setDescription(event_.target.value); }}
                    placeholder="Краткое описание идеи..."
                    disabled={isSubmitting}
                />
            </div>
            <div className={styles.field}>
                <label className={styles.label}>Содержимое</label>
                <MdEditor
                    value={content}
                    onChange={setContent}
                    placeholder="Напишите подробное содержимое идеи..."
                    height={400}
                    disabled={isSubmitting}
                />
            </div>

            {error !== undefined && <p className={styles.error}>{error}</p>}

            <div className={styles.actions}>
                <Button type="button" onClick={onCancel} disabled={isSubmitting}>
                    Отмена
                </Button>
                <Button type="submit" disabled={isSubmitting || !title}>
                    Сохранить изменения
                </Button>
            </div>
        </form>
    );
});
