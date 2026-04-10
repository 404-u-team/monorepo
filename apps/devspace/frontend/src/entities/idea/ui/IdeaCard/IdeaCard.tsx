import { Link } from "@tanstack/react-router";
import { clsx } from "clsx";
import { Heart, Eye } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useEffect, useState, type JSX } from "react";

import { fetchUserById, type IUserResponse } from "@/entities/user";
import { useStore } from "@/shared/lib/store";
import { Badge, IconCounter } from "@/shared/ui";

import { fetchIdeaById, toggleIdeaFavorite } from "../../api/ideaApi";
import type { IIdea } from "../../model/IIdea";
import { IdeaCardSkeleton } from "../IdeaCardSkeleton/IdeaCardSkeleton";

import styles from "./IdeaCard.module.scss";

export interface IdeaCardProps {
  ideaId: string;
  href?: string | undefined;
  fromRoute?: string | undefined;
  className?: string | undefined;
}

export const IdeaCard = observer(function IdeaCard({
  ideaId,
  href,
  fromRoute,
  className,
}: IdeaCardProps): JSX.Element {
  const { userStore } = useStore();

  const [idea, setIdea] = useState<IIdea | undefined>(undefined);
  const [author, setAuthor] = useState<IUserResponse | undefined>(undefined);
  const [isFavorite, setIsFavorite] = useState(false);
  const [favoritesCount, setFavoritesCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function load(): Promise<void> {
      try {
        const ideaData = await fetchIdeaById(ideaId);
        if (cancelled) return;
        setIdea(ideaData);
        setIsFavorite(ideaData.is_favorite ?? false);
        setFavoritesCount(ideaData.favorites_count);

        const authorData = await fetchUserById(ideaData.author_id);
        // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition -- cancelled may change across await
        if (cancelled) return;
        setAuthor(authorData);
      } catch {
        // Errors can be handled with a future error state
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    void load();
    return (): void => {
      cancelled = true;
    };
  }, [ideaId]);

  const handleFavoriteClick = async (): Promise<void> => {
    if (!userStore.isAuthenticated) return;
    try {
      const result = await toggleIdeaFavorite(ideaId);
      setIsFavorite(result.is_favorite);
      setFavoritesCount((previous) => (result.is_favorite ? previous + 1 : previous - 1));
    } catch {
      // Silently fail for now
    }
  };

  if (isLoading || idea === undefined) {
    return <IdeaCardSkeleton className={className} />;
  }

  const targetHref = href ?? (ideaId !== "" ? `/idea/${ideaId}` : undefined);

  const cardContent = (
    <>
      <div className={styles.imageWrapper}>
        <div className={styles.imagePlaceholder} />
        {idea.category !== undefined && idea.category !== "" && (
          <Badge className={styles.categoryBadge}>{idea.category}</Badge>
        )}
      </div>

      <div className={styles.content}>
        <h3 className={styles.title}>{idea.title}</h3>
        <p className={styles.description}>{idea.description}</p>
      </div>

      <div className={styles.footer}>
        <div className={styles.author}>
          {author !== undefined && (
            <>
              <img
                className={styles.avatar}
                src={author.avatar_url}
                alt={author.nickname}
                onError={(event) => {
                  (event.target as HTMLImageElement).style.display = "none";
                }}
              />
              <span className={styles.nickname}>{author.nickname}</span>
            </>
          )}
        </div>

        <div className={styles.stats}>
          <IconCounter
            icon={<Heart size={16} fill={isFavorite ? "currentColor" : "none"} />}
            count={favoritesCount}
            active={isFavorite}
            onClick={(event) => {
              event.preventDefault();
              event.stopPropagation();
              void handleFavoriteClick();
            }}
          />
          <IconCounter icon={<Eye size={16} />} count={idea.views_count} />
        </div>
      </div>
    </>
  );

  if (targetHref !== undefined) {
    return (
      <Link
        to={targetHref}
        {...(fromRoute !== undefined
          ? { state: (previous) => ({ ...previous, backTo: fromRoute }) }
          : {})}
        className={clsx(styles.card, styles.link, className)}
      >
        {cardContent}
      </Link>
    );
  }

  return <article className={clsx(styles.card, className)}>{cardContent}</article>;
});
