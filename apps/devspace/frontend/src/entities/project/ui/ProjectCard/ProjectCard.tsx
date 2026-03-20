import { useEffect, useState, type JSX } from "react";
import { Heart } from "lucide-react";
import { clsx } from "clsx";
import { Badge } from "@/shared/ui";
import { fetchUserById } from "@/entities/user";
import type { IProject } from "../../model/IProject";
import type { IProjectSlotSkill } from "../../model/IProjectSlot";
import {
  fetchProjectById,
  type IProjectDetailResponse,
} from "../../api/projectApi";
import { ProjectCardSkeleton } from "../ProjectCardSkeleton/ProjectCardSkeleton";
import styles from "./ProjectCard.module.scss";

export interface ProjectCardProps {
  projectId: string;
  to?: string | undefined;
  className?: string | undefined;
}

interface SlotUser {
  id: string;
  avatar_uri: string;
  nickname: string;
}

const STATUS_LABEL: Record<IProject["status"], string> = {
  open: "In Progress",
  closed: "Closed",
};

export function ProjectCard({ projectId, to, className }: ProjectCardProps): JSX.Element {
    const [project, setProject] = useState<IProjectDetailResponse | undefined>(undefined);
    const [slotUsers, setSlotUsers] = useState<SlotUser[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        let cancelled = false;

        async function load(): Promise<void> {
            try {
                const data = await fetchProjectById(projectId);
                if (cancelled) return;
                setProject(data);

                const occupiedUserIds = (data.slots ?? [])
                    .map((s) => s.user_id)
                    .filter((uid): uid is string => uid !== null);

                const users = await Promise.all(
                    occupiedUserIds.map((uid) => fetchUserById(uid)),
                );
                // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition -- cancelled may change across await
                if (cancelled) return;
                setSlotUsers(users.map((u) => ({ id: u.id, avatar_uri: u.avatar_uri, nickname: u.nickname })));
            } catch {
                // handled by future error state
            } finally {
                if (!cancelled) setIsLoading(false);
            }
        }

        void load();
        return (): void => { cancelled = true; };
    }, [projectId]);

    if (isLoading || project === undefined) {
        return <ProjectCardSkeleton className={className} />;
    }

    const openSlots = (project.slots ?? []).filter((s) => s.status === 'open');

    // Collect unique primary (1st level) skills only — secondary skills are shown on detail page
    // TODO: When user primary skill is implemented, show slot skills relevant to the user's primary skill
    const uniqueSkills: IProjectSlotSkill[] = [];
    for (const slot of (project.slots ?? [])) {
        if (slot.skill !== undefined && !uniqueSkills.some((s) => s.id === slot.skill?.id)) {
            uniqueSkills.push(slot.skill);
        }
    }
  }

  const Wrapper = to !== undefined ? "a" : "article";
  const wrapperProps = to !== undefined ? { href: to } : {};

  return (
    <Wrapper
      {...wrapperProps}
      className={clsx(styles.card, to !== undefined && styles.link, className)}
    >
      <div className={styles.imageWrapper}>
        <div className={styles.imagePlaceholder} />
        <Badge className={clsx(styles.statusBadge, styles[project.status])}>
          {STATUS_LABEL[project.status]}
        </Badge>
      </div>

      <div className={styles.body}>
        <div className={styles.header}>
          <h3 className={styles.title}>{project.title}</h3>
          <button
            type="button"
            className={styles.favoriteBtn}
            aria-label="Favorite"
          >
            <Heart size={20} />
          </button>
        </div>

        <p className={styles.description}>{project.description}</p>

        {uniqueSkills.length > 0 && (
          <div className={styles.skills}>
            {uniqueSkills.map((skill) => (
              <Badge key={skill.id} color={skill.color}>
                {skill.name}
              </Badge>
            ))}
          </div>
        )}

        {openSlots.length > 0 && (
          <div className={styles.openSlots}>
            <span className={styles.openSlotsLabel}>Свободные слоты:</span>
            <div className={styles.openSlotsList}>
              {openSlots.map((slot) => (
                <Badge key={slot.id} color={slot.skill?.color}>
                  {slot.skill?.name ?? slot.title}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {slotUsers.length > 0 && (
          <div className={styles.members}>
            <div className={styles.avatars}>
              {slotUsers.slice(0, 3).map((user) => (
                <img
                  key={user.id}
                  className={styles.avatar}
                  src={user.avatar_uri}
                  alt={user.nickname}
                  onError={(event) => {
                    (event.target as HTMLImageElement).style.display = "none";
                  }}
                />
              ))}
            </div>
            {slotUsers.length > 3 && (
              <span className={styles.membersMore}>
                +{slotUsers.length - 3}
              </span>
            )}
          </div>
        )}
      </div>
    </Wrapper>
  );
}
