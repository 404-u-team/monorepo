import { type JSX } from 'react';
import { User } from 'lucide-react';
import { Badge, Skeleton } from '@/shared/ui';
import type { IUserResponse } from '@/entities/user';
import styles from './UserPublicProfile.module.scss';

export interface UserPublicProfileProps {
    user: IUserResponse;
}

export function UserPublicProfile({ user }: UserPublicProfileProps): JSX.Element {
    return (
        <div className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.avatarWrapper}>
                    {user.avatar_uri ? (
                        <img
                            src={user.avatar_uri}
                            alt={user.nickname}
                            className={styles.avatar}
                            onError={(event) => {
                                (event.target as HTMLImageElement).style.display = 'none'
                            }}
                        />
                    ) : (
                        <div className={styles.avatarPlaceholder}>
                            <User size={48} />
                        </div>
                    )}
                </div>

                <div className={styles.info}>
                    <h1 className={styles.nickname}>{user.nickname}</h1>
                    {user.main_role && (
                        <p className={styles.mainRole}>{user.main_role}</p>
                    )}
                </div>
            </div>

            {user.bio && (
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>О себе</h2>
                    <p className={styles.bio}>{user.bio}</p>
                </section>
            )}

            {user.skills.length > 0 && (
                <section className={styles.section}>
                    <h2 className={styles.sectionTitle}>Навыки</h2>
                    <div className={styles.skills}>
                        {user.skills.map((skill) => (
                            <Badge key={skill.id}>{skill.name}</Badge>
                        ))}
                    </div>
                </section>
            )}
        </div>
    );
}

export function UserPublicProfileSkeleton(): JSX.Element {
    return (
        <div className={styles.wrapper}>
            <div className={styles.header}>
                <Skeleton width={120} height={120} borderRadius={60} />
                <div className={styles.info}>
                    <Skeleton width={200} height={32} />
                    <Skeleton width={140} height={20} />
                </div>
            </div>

            <div className={styles.section}>
                <Skeleton width={80} height={24} />
                <div style={{ marginTop: '12px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    <Skeleton width="100%" height={16} />
                    <Skeleton width="90%" height={16} />
                    <Skeleton width="70%" height={16} />
                </div>
            </div>

            <div className={styles.section}>
                <Skeleton width={80} height={24} />
                <div style={{ marginTop: '12px', display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
                    <Skeleton width={80} height={28} borderRadius={100} />
                    <Skeleton width={70} height={28} borderRadius={100} />
                    <Skeleton width={90} height={28} borderRadius={100} />
                    <Skeleton width={65} height={28} borderRadius={100} />
                </div>
            </div>
        </div>
    );
}
