import { useState, type JSX } from 'react';
import { observer } from 'mobx-react-lite';
import { Button, Input, Dropdown } from '@/shared/ui';
import { updateProject, type IProject } from '@/entities/project';
import styles from './EditProjectForm.module.scss';

interface EditProjectFormProps {
    project: IProject;
    onSuccess: (updated: IProject) => void;
    onCancel: () => void;
}

const statusOptions = [
    { label: 'Открытый', value: 'open' },
    { label: 'Закрытый', value: 'closed' },
];

export const EditProjectForm = observer(function EditProjectForm({
    project,
    onSuccess,
    onCancel,
}: EditProjectFormProps): JSX.Element {
    const [title, setTitle] = useState(project.title);
    const [description, setDescription] = useState(project.description);
    const [status, setStatus] = useState<'open' | 'closed'>(project.status);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | undefined>(undefined);

    const handleSubmit = async (event: React.SyntheticEvent): Promise<void> => {
        event.preventDefault();
        if (!title) return;

        setIsSubmitting(true);
        setError(undefined);
        try {
            const updated = await updateProject(project.id, { title, description, status });
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
                <label className={styles.label} htmlFor="project-title">Название</label>
                <Input
                    id="project-title"
                    value={title}
                    onChange={(event_) => { setTitle(event_.target.value); }}
                    placeholder="Введите название..."
                    required
                    disabled={isSubmitting}
                />
            </div>
            <div className={styles.field}>
                <label className={styles.label} htmlFor="project-description">Описание</label>
                <textarea
                    id="project-description"
                    className={styles.textarea}
                    value={description}
                    onChange={(event_) => { setDescription(event_.target.value); }}
                    placeholder="Опишите цели и задачи проекта..."
                    disabled={isSubmitting}
                />
            </div>
            <div className={styles.field}>
                <label className={styles.label}>Статус</label>
                <Dropdown
                    options={statusOptions}
                    value={status}
                    onChange={(value) => { setStatus(value as 'open' | 'closed'); }}
                />
            </div>

            {error !== undefined && <p className={styles.error}>{error}</p>}

            <div className={styles.actions}>
                <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting}>
                    Отмена
                </Button>
                <Button type="submit" disabled={isSubmitting || !title}>
                    Сохранить изменения
                </Button>
            </div>
        </form>
    );
});
