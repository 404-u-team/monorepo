import { useEffect, useState, type JSX } from 'react';
import { useParams, Link } from '@tanstack/react-router';
import { observer } from 'mobx-react-lite';
import { Heart, Eye, ArrowLeft, Rocket } from 'lucide-react';
import { useStore } from '@/shared/lib/store';
import { Button, Badge, IconCounter, Skeleton } from '@/shared/ui';
import { fetchUserById, type IUserResponse } from '@/entities/user';
import { fetchIdeaById, toggleIdeaFavorite, createProjectFromIdea, type IIdea } from '@/entities/idea';
import styles from './IdeaDetail.module.scss';

export const IdeaDetail = observer(function IdeaDetail(): JSX.Element {
    const { ideaId } = useParams({ from: '/idea/$ideaId' });
    const { userStore } = useStore();

    const [idea, setIdea] = useState<IIdea | undefined>(undefined);
    const [author, setAuthor] = useState<IUserResponse | undefined>(undefined);
    const [isFavorite, setIsFavorite] = useState(false);
    const [favoritesCount, setFavoritesCount] = useState(0);
    const [isLoading, setIsLoading] = useState(true);
    const [isCreatingProject, setIsCreatingProject] = useState(false);

    useEffect(() => {
        let cancelled = false;

        async function load(): Promise<void> {
            try {
                const ideaData = await fetchIdeaById(ideaId);
                if (cancelled) return;
                setIdea(ideaData);
                setFavoritesCount(ideaData.favorites_count);

                const authorData = await fetchUserById(ideaData.author_id);
                if (cancelled) return;
                setAuthor(authorData);
            } catch {
                // handle error
            } finally {
                if (!cancelled) {
                    setIsLoading(false);
                }
            }
        }

        void load();
        return (): void => { cancelled = true; };
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

    if (isLoading || idea === undefined) {
        return (
            <div className={styles.container}>
                <Skeleton className={styles.skeletonTitle} />
                <Skeleton className={styles.skeletonText} />
                <Skeleton className={styles.skeletonText} />
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
                    <Badge>{idea.category}</Badge>
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
                            onClick={handleFavoriteClick}
                        />
                        <IconCounter
                            icon={<Eye size={20} />}
                            count={idea.views_count}
                        />
                    </div>
                    {userStore.isAuthenticated && (
                        <Button
                            onClick={handleCreateProject}
                            disabled={isCreatingProject}
                            className={styles.createProjectBtn}
                        >
                            <Rocket size={18} />
                            Создать проект
                        </Button>
                    )}
                </div>
            </header>

            <div className={styles.content}>
                <div className={styles.descriptionSection}>
                    <h2 className={styles.sectionTitle}>Описание идеи</h2>
                    <p className={styles.descriptionText}>{idea.description}</p>
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
