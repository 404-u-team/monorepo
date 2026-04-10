import { useParams, useNavigate, useLocation, useRouter } from "@tanstack/react-router";
import { Heart, Eye, ArrowLeft, Rocket, Pencil, Trash2 } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useEffect, useState, useRef, type JSX } from "react";

import { fetchIdeaById, toggleIdeaFavorite, deleteIdea, type IIdea } from "@/entities/idea";
import { fetchUserById, type IUserResponse, type UserStore } from "@/entities/user";
import { EditIdeaForm } from "@/features/idea/edit";
import { useStore } from "@/shared/lib/store";
import { Button, Badge, IconCounter, Skeleton, MdRenderer, ConfirmModal } from "@/shared/ui";

import styles from "./IdeaDetail.module.scss";

export const IdeaDetail = observer(function IdeaDetail(): JSX.Element {
  const { ideaId } = useParams({ from: "/idea/$ideaId" });
  const { userStore } = useStore() as unknown as { userStore: UserStore };
  const navigate = useNavigate();
  const location = useLocation();
  const router = useRouter();
  const backTo = (location.state as { backTo?: string } | null)?.backTo;

  const [idea, setIdea] = useState<IIdea | undefined>(undefined);
  const [author, setAuthor] = useState<IUserResponse | undefined>(undefined);
  const [isFavorite, setIsFavorite] = useState(false);
  const [favoritesCount, setFavoritesCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const isCancelled = useRef(false);

  useEffect(() => {
    isCancelled.current = false;
    setIsLoading(true);
    setIdea(undefined);
    setAuthor(undefined);
    setFavoritesCount(0);
    setIsFavorite(false);

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
    return (): void => {
      isCancelled.current = true;
    };
  }, [ideaId]);

  const handleFavoriteClick = async (): Promise<void> => {
    if (!userStore.isAuthenticated) return;
    try {
      const result = await toggleIdeaFavorite(ideaId);
      setIsFavorite(result.is_favorite);
      setFavoritesCount((previous) => (result.is_favorite ? previous + 1 : previous - 1));
    } catch {
      // handle error
    }
  };

  const handleCreateProject = (): void => {
    if (!idea) return;
    void navigate({
      to: "/project/new",
      search: {
        title: idea.title,
        description: idea.description,
        idea_id: idea.id,
      },
    });
  };

  const handleDelete = async (): Promise<void> => {
    setIsDeleting(true);
    try {
      await deleteIdea(ideaId);
      void navigate({ to: "/ideas" });
    } catch {
      setIsDeleting(false);
    } finally {
      setIsDeleteModalOpen(false);
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

  const isAuthor = userStore.isAuthenticated && userStore.user?.id === idea.author_id;

  if (isEditing) {
    return (
      <div className={styles.container}>
        <button
          type="button"
          className={styles.backLink}
          onClick={() => {
            setIsEditing(false);
          }}
        >
          <ArrowLeft size={16} />
          Назад к идее
        </button>
        <h1 className={styles.editTitle}>Редактировать идею</h1>
        <EditIdeaForm
          idea={idea}
          onSuccess={handleEditSuccess}
          onCancel={() => {
            setIsEditing(false);
          }}
        />
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {backTo !== undefined && (
        <button
          type="button"
          className={styles.backLink}
          onClick={() => {
            router.history.back();
          }}
        >
          <ArrowLeft size={16} />
          Назад к списку
        </button>
      )}

      <header className={styles.header}>
        <div className={styles.titleWrapper}>
          <h1 className={styles.title}>{idea.title}</h1>
          {idea.category !== undefined && idea.category !== "" && <Badge>{idea.category}</Badge>}
        </div>

        <div className={styles.actions}>
          <div className={styles.stats}>
            <IconCounter
              icon={<Heart size={20} fill={isFavorite ? "currentColor" : "none"} />}
              count={favoritesCount}
              active={isFavorite}
              onClick={() => {
                void handleFavoriteClick();
              }}
            />
            <IconCounter icon={<Eye size={20} />} count={idea.views_count} />
          </div>
          <div className={styles.actionButtons}>
            {isAuthor && (
              <>
                <Button
                  onClick={() => {
                    setIsEditing(true);
                  }}
                >
                  <Pencil size={18} />
                  Редактировать
                </Button>
                <Button
                  variant="outline"
                  onClick={() => {
                    setIsDeleteModalOpen(true);
                  }}
                  disabled={isDeleting}
                >
                  <Trash2 size={18} />
                  Удалить
                </Button>
              </>
            )}
            {userStore.isAuthenticated && (
              <Button onClick={handleCreateProject} className={styles.createProjectBtn}>
                <Rocket size={18} />
                Создать проект
              </Button>
            )}
          </div>
        </div>
      </header>

      <div className={styles.content}>
        <div className={styles.mainSection}>
          {idea.description !== "" && (
            <div className={styles.descriptionSection}>
              <h2 className={styles.sectionTitle}>Описание идеи</h2>
              <p className={styles.descriptionText}>{idea.description}</p>
            </div>
          )}

          {idea.content !== undefined && idea.content !== "" && (
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
                  src={author.avatar_url}
                  alt={author.nickname}
                  onError={(event) => {
                    const target = event.currentTarget;
                    target.onerror = undefined as unknown as typeof target.onerror;
                    target.style.display = "none";
                  }}
                />
                <span className={styles.nickname}>{author.nickname}</span>
              </div>
            )}
          </div>
        </aside>
      </div>
      <ConfirmModal
        isOpen={isDeleteModalOpen}
        title="Удалить идею?"
        description="Вы уверены, что хотите удалить эту идею? Это действие нельзя отменить."
        severity="danger"
        confirmLabel="Удалить"
        cancelLabel="Отмена"
        isLoading={isDeleting}
        onConfirm={() => {
          void handleDelete();
        }}
        onCancel={() => {
          setIsDeleteModalOpen(false);
        }}
      />
    </div>
  );
});
