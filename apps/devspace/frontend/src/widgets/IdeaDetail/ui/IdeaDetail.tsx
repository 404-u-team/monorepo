import { useEffect, useState, useRef, type JSX } from 'react';
import { useParams, Link } from '@tanstack/react-router';
import { observer } from 'mobx-react-lite';
import { Heart, Eye, ArrowLeft, Rocket, Pencil } from 'lucide-react';
import { useStore } from '@/shared/lib/store';
import { Button, Badge, IconCounter, Skeleton, MdRenderer } from '@/shared/ui';
import { fetchUserById, type IUserResponse, type UserStore } from '@/entities/user';
import { fetchIdeaById, toggleIdeaFavorite, createProjectFromIdea, type IIdea } from '@/entities/idea';
import { EditIdeaForm } from '@/features/idea/edit';
import styles from './IdeaDetail.module.scss';

export const IdeaDetail = observer(function IdeaDetail(): JSX.Element {
    const { ideaId } = useParams({ from: '/idea/$ideaId' });
    const { userStore } = useStore() as unknown as { userStore: UserStore };

    const [idea, setIdea] = useState<IIdea | undefined>(undefined);
    const [author, setAuthor] = useState<IUserResponse | undefined>(undefined);
    const [isFavorite, setIsFavorite] = useState(false);
    const [favoritesCount, setFavoritesCount] = useState(0);
    const [isLoading, setIsLoading] = useState(true);
    const [isCreatingProject, setIsCreatingProject] = useState(false);
    const [isEditing, setIsEditing] = useState(false);

    const isCancelled = useRef(false);

    useEffect(() => {
        isCancelled.current = false;

        async function load(): Promise<void> {
            try {
                const ideaData = await fetchIdeaById(ideaId);
                if (isCancelled.current) return;
                setIdea(ideaData);
                setFavoritesCount(ideaData.favorites_count);

                const authorData = await fetchUserById(ideaData.author_id);
                // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition -- isCancelled.current may change after await
                if (isCancelled.current) return;
                setAuthor(authorData);
            } catch {
                // handle error
            } finally {
                if (!isCancelled.current) {
                    setIsLoading(false);
                }
            }
        }

        void load();
        return (): void => { isCancelled.current = true; };
    }, [ideaId]);

    const handleFavoriteClick = async (): Promise<void> => {
        if (!userStore.isAuthenticated) return;
        try {
            const result = await toggleIdeaFavorite(ideaId);
            setIsFavorite(result.is_favorite);
            setFavoritesCount((previous) =>
                result.is_favorite ? previous + 1 : previous - 1,
            );
        } catch {
            // handle error
        }
    };

    const handleCreateProject = async (): Promise<void> => {
        if (!userStore.isAuthenticated) return;
        setIsCreatingProject(true);
        try {
            await createProjectFromIdea(ideaId);
            // Redirect to projects or show success
            alert('Проект успешно создан!');
        } catch {
            // handle error
        } finally {
            setIsCreatingProject(false);
        }
    };

    const handleEditSuccess = (updated: IIdea): void => {
        setIdea(updated);
        setIsEditing(false);
    };

    if (isLoading || idea === undefined) {
        return (
            <div className={styles.container}>
                <Skeleton className={styles.skeletonTitle} />
                <Skeleton className={styles.skeletonText} />
                <Skeleton className={styles.skeletonText} />
            </div>
        );
    }

    const isAuthor =
        userStore.isAuthenticated &&
        userStore.user?.id === idea.author_id;

    if (isEditing) {
        return (
            <div className={styles.container}>
                <button type="button" className={styles.backLink} onClick={() => { setIsEditing(false); }}>
                    <ArrowLeft size={16} />
                    Назад к идее
                </button>
                <h1 className={styles.editTitle}>Редактировать идею</h1>
                <EditIdeaForm
                    idea={idea}
                    onSuccess={handleEditSuccess}
                    onCancel={() => { setIsEditing(false); }}
                />
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <Link to="/ideas" className={styles.backLink}>
                <ArrowLeft size={16} />
                Назад к списку
            </Link>

            <header className={styles.header}>
                <div className={styles.titleWrapper}>
                    <h1 className={styles.title}>{idea.title}</h1>
                    {idea.category !== undefined && idea.category !== '' && (
                        <Badge>{idea.category}</Badge>
                    )}
                </div>

                <div className={styles.actions}>
                    <div className={styles.stats}>
                        <IconCounter
                            icon={
                                <Heart
                                    size={20}
                                    fill={isFavorite ? 'currentColor' : 'none'}
                                />
                            }
                            count={favoritesCount}
                            active={isFavorite}
                            onClick={() => { void handleFavoriteClick(); }}
                        />
                        <IconCounter
                            icon={<Eye size={20} />}
                            count={idea.views_count}
                        />
                    </div>
                    <div className={styles.actionButtons}>
                        {isAuthor && (
                            <Button onClick={() => { setIsEditing(true); }}>
                                <Pencil size={18} />
                                Редактировать
                            </Button>
                        )}
                        {userStore.isAuthenticated && (
                            <Button
                                onClick={() => { void handleCreateProject(); }}
                                disabled={isCreatingProject}
                                className={styles.createProjectBtn}
                            >
                                <Rocket size={18} />
                                Создать проект
                            </Button>
                        )}
                    </div>
                </div>
            </header>

            <div className={styles.content}>
                <div className={styles.mainSection}>
                    {idea.description !== '' && (
                        <div className={styles.descriptionSection}>
                            <h2 className={styles.sectionTitle}>Описание идеи</h2>
                            <p className={styles.descriptionText}>{idea.description}</p>
                        </div>
                    )}

                    {idea.content !== undefined && idea.content !== '' && (
                        <div className={styles.contentSection}>
                            <h2 className={styles.sectionTitle}>Содержимое</h2>
                            <MdRenderer source={idea.content} className={styles.markdownContent} />
                        </div>
                    )}
                </div>

                <aside className={styles.sidebar}>
                    <div className={styles.authorCard}>
                        <h3 className={styles.sidebarTitle}>Автор идеи</h3>
                        {author !== undefined && (
                            <div className={styles.authorInfo}>
                                <img
                                    className={styles.avatar}
                                    src={author.avatar_uri}
                                    alt={author.nickname}
                                />
                                <span className={styles.nickname}>{author.nickname}</span>
                            </div>
                        )}
                    </div>
                </aside>
            </div>
        </div>
    );
});
