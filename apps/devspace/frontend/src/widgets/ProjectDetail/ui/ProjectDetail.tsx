import { useEffect, useState, useRef, useCallback, type JSX } from 'react';
import { useParams, Link, useNavigate } from '@tanstack/react-router';
import { observer } from 'mobx-react-lite';
import { ArrowLeft, Pencil, Trash2, Plus, X, Check, UserCheck } from 'lucide-react';
import { useStore } from '@/shared/lib/store';
import { Button, Badge, Skeleton, ConfirmModal, SkillSearch, SkillMultiSelect, UserAvatar, MdRenderer, type SkillSearchOption, type SkillMultiSelectOption } from '@/shared/ui';
import { fetchUserById, type IUserResponse, type UserStore } from '@/entities/user';
import {
    fetchProjectById,
    fetchProjectSlots,
    deleteProject,
    createProjectSlot,
    deleteProjectSlot,
    applyToSlot,
    getProjectRequests,
    acceptRequest,
    rejectRequest,
    type IProjectDetailResponse,
    type IProjectSlot,
    type IRequest,
} from '@/entities/project';
import { fetchSkills } from '@/entities/skill';
import { EditProjectForm } from '@/features/project/edit';
import styles from './ProjectDetail.module.scss';

function slotStatusLabel(isOccupied: boolean, status: 'open' | 'closed'): string {
    if (isOccupied) return 'Занято';
    if (status === 'open') return 'Открыт';
    return 'Закрыт';
}

export const ProjectDetail = observer(function ProjectDetail(): JSX.Element {
    const { projectId } = useParams({ from: '/project/$projectId' });
    const { userStore } = useStore() as unknown as { userStore: UserStore };
    const navigate = useNavigate();

    const [project, setProject] = useState<IProjectDetailResponse | undefined>(undefined);
    const [leader, setLeader] = useState<IUserResponse | undefined>(undefined);
    const [requests, setRequests] = useState<IRequest[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [isEditing, setIsEditing] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    // Slot creation state
    const [isCreatingSlot, setIsCreatingSlot] = useState(false);
    const [newSlotTitle, setNewSlotTitle] = useState('');
    const [newSlotDescription, setNewSlotDescription] = useState('');
    const [newSlotPrimarySkills, setNewSlotPrimarySkills] = useState<SkillMultiSelectOption[]>([]);
    const [newSlotSecondarySkills, setNewSlotSecondarySkills] = useState<SkillMultiSelectOption[]>([]);
    const [slotCreateError, setSlotCreateError] = useState<string | undefined>(undefined);
    const [isSubmittingSlot, setIsSubmittingSlot] = useState(false);

    // Skill dropdown filter on detail page
    const [selectedSkillFilter, setSelectedSkillFilter] = useState<SkillSearchOption | undefined>(undefined);

    // Apply state per slot
    const [applyingSlotId, setApplyingSlotId] = useState<string | undefined>(undefined);

    // Delete confirmation modals
    const [isProjectDeleteModalOpen, setIsProjectDeleteModalOpen] = useState(false);
    const [slotIdToDelete, setSlotIdToDelete] = useState<string | undefined>(undefined);

    const isCancelled = useRef(false);

    const loadProject = async (): Promise<void> => {
        try {
            const [projectData, slotsData] = await Promise.all([
                fetchProjectById(projectId),
                fetchProjectSlots(projectId),
            ]);
            if (isCancelled.current) return;
            setProject({ ...projectData, slots: slotsData });

            const leaderData = await fetchUserById(projectData.leader_id);
            // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition -- isCancelled.current may change after await
            if (isCancelled.current) return;
            setLeader(leaderData);
        } catch {
            // handle error
        } finally {
            if (!isCancelled.current) {
                setIsLoading(false);
            }
        }
    };

    const loadRequests = async (): Promise<void> => {
        if (!userStore.isAuthenticated) return;
        try {
            const requestData = await getProjectRequests(projectId);
            if (!isCancelled.current) {
                setRequests(requestData);
            }
        } catch {
            // not leader — no access, ignore
        }
    };

    useEffect(() => {
        isCancelled.current = false;
        setIsLoading(true);
        setProject(undefined);
        setLeader(undefined);
        setRequests([]);

        void loadProject();
        return (): void => { isCancelled.current = true; };
    }, [projectId]); // eslint-disable-line react-hooks/exhaustive-deps

    useEffect(() => {
        if (project !== undefined && userStore.isAuthenticated && userStore.user?.id === project.leader_id) {
            void loadRequests();
        }
    }, [project, userStore.isAuthenticated, userStore.user?.id]); // eslint-disable-line react-hooks/exhaustive-deps

    const handleDelete = async (): Promise<void> => {
        setIsDeleting(true);
        try {
            await deleteProject(projectId);
            void navigate({ to: '/projects' });
        } catch {
            setIsDeleting(false);
        } finally {
            setIsProjectDeleteModalOpen(false);
        }
    };

    // Loaders for skill search components
    const loadPrimarySkills = useCallback(async (query: string): Promise<SkillSearchOption[]> => {
        const skills = await fetchSkills({ search: query || undefined, limit: 20 });
        // Primary skills = root level (parent_id is null)
        return skills
            .filter((skill) => skill.parent_id === null)
            .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
    }, []);

    const loadSecondarySkills = useCallback(async (query: string): Promise<SkillMultiSelectOption[]> => {
        const skills = await fetchSkills({ search: query || undefined, limit: 30 });
        // Secondary skills = children (parent_id !== null)
        // Prioritize those linked to selected primary skills
        const secondary = skills.filter((skill) => skill.parent_id !== null);
        if (newSlotPrimarySkills.length > 0) {
            const primaryIds = new Set(newSlotPrimarySkills.map((s) => s.id));
            const linked = secondary.filter((skill) => skill.parent_id !== null && primaryIds.has(skill.parent_id));
            const other = secondary.filter((skill) => skill.parent_id === null || !primaryIds.has(skill.parent_id));
            return [...linked, ...other].map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
        }
        return secondary.map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
    }, [newSlotPrimarySkills]);

    const handleCreateSlot = async (): Promise<void> => {
        if (!newSlotTitle || newSlotPrimarySkills.length === 0) {
            setSlotCreateError('Заполните название и выберите хотя бы один основной навык');
            return;
        }
        setIsSubmittingSlot(true);
        setSlotCreateError(undefined);
        try {
            const payload: Parameters<typeof createProjectSlot>[1] = {
                title: newSlotTitle,
                description: newSlotDescription,
                primary_skills_id: newSlotPrimarySkills.map((skill) => skill.id),
                status: 'open',
            };
            if (newSlotSecondarySkills.length > 0) {
                payload.secondary_skills_id = newSlotSecondarySkills.map((skill) => skill.id);
            }
            const slot = await createProjectSlot(projectId, payload);
            setProject((previous) => {
                if (!previous) return previous;
                return { ...previous, slots: [...(previous.slots ?? []), slot] };
            });
            setIsCreatingSlot(false);
            setNewSlotTitle('');
            setNewSlotDescription('');
            setNewSlotPrimarySkills([]);
            setNewSlotSecondarySkills([]);
        } catch {
            setSlotCreateError('Ошибка при создании слота');
        } finally {
            setIsSubmittingSlot(false);
        }
    };

    const handleDeleteSlot = async (slotId: string): Promise<void> => {
        try {
            await deleteProjectSlot(projectId, slotId);
            setProject((previous) => {
                if (!previous) return previous;
                return { ...previous, slots: (previous.slots ?? []).filter((slot) => slot.id !== slotId) };
            });
            setRequests((previous) => previous.filter((request) => request.slot_id !== slotId));
        } catch {
            // handle error
        } finally {
            setSlotIdToDelete(undefined);
        }
    };

    const handleApply = async (slotId: string): Promise<void> => {
        setApplyingSlotId(slotId);
        try {
            await applyToSlot(projectId, slotId);
            // Refresh requests if leader, or just show feedback
        } catch {
            // handle error
        } finally {
            setApplyingSlotId(undefined);
        }
    };

    const handleAcceptRequest = async (requestId: string): Promise<void> => {
        try {
            const updated = await acceptRequest(requestId);
            setRequests((previous) =>
                previous.map((request) => request.id === requestId ? updated : request),
            );
            // Update slot to show it's occupied
            if (project) {
                const relatedRequest = requests.find((r) => r.id === requestId);
                if (relatedRequest) {
                    setProject((previous) => {
                        if (!previous) return previous;
                        return {
                            ...previous,
                            slots: (previous.slots ?? []).map((slot) =>
                                slot.id === relatedRequest.slot_id
                                    ? { ...slot, user_id: relatedRequest.user_id }
                                    : slot,
                            ),
                        };
                    });
                }
            }
        } catch {
            // handle error
        }
    };

    const handleRejectRequest = async (requestId: string): Promise<void> => {
        try {
            const updated = await rejectRequest(requestId);
            setRequests((previous) =>
                previous.map((request) => request.id === requestId ? updated : request),
            );
        } catch {
            // handle error
        }
    };

    if (isLoading || project === undefined) {
        return (
            <div className={styles.container}>
                <Skeleton className={styles.skeletonTitle} />
                <Skeleton className={styles.skeletonText} />
                <Skeleton className={styles.skeletonText} />
            </div>
        );
    }

    const isLeader = userStore.isAuthenticated && userStore.user?.id === project.leader_id;

    const pendingRequestsBySlot = (slotId: string): IRequest[] =>
        requests.filter((request) => request.slot_id === slotId && request.status === 'pending');

    if (isEditing) {
        return (
            <div className={styles.container}>
                <button type="button" className={styles.backLink} onClick={() => { setIsEditing(false); }}>
                    <ArrowLeft size={16} />
                    Назад к проекту
                </button>
                <h1 className={styles.editTitle}>Редактировать проект</h1>
                <EditProjectForm
                    project={project}
                    onSuccess={(updated) => {
                        setProject((previous) => previous ? { ...previous, ...updated } : previous);
                        setIsEditing(false);
                    }}
                    onCancel={() => { setIsEditing(false); }}
                />
            </div>
        );
    }

    return (
        <div className={styles.container}>
            <Link to="/projects" className={styles.backLink}>
                <ArrowLeft size={16} />
                Назад к списку
            </Link>

            <header className={styles.header}>
                <div className={styles.titleWrapper}>
                    <h1 className={styles.title}>{project.title}</h1>
                    <Badge color={project.status === 'open' ? '039855' : undefined}>
                        {project.status === 'open' ? 'Открытый' : 'Закрытый'}
                    </Badge>
                </div>

                {isLeader && (
                    <div className={styles.actionButtons}>
                        <Button onClick={() => { setIsEditing(true); }}>
                            <Pencil size={18} />
                            Редактировать
                        </Button>
                        <Button
                            variant="outline"
                            onClick={() => { setIsProjectDeleteModalOpen(true); }}
                            disabled={isDeleting}
                        >
                            <Trash2 size={18} />
                            Удалить
                        </Button>
                    </div>
                )}
            </header>

            <div className={styles.content}>
                <div className={styles.mainSection}>
                    {project.description !== '' && (
                        <div className={styles.descriptionSection}>
                            <h2 className={styles.sectionTitle}>Описание</h2>
                            <p className={styles.descriptionText}>{project.description}</p>
                        </div>
                    )}

                    {project.content !== null && project.content !== '' && (
                        <div className={styles.descriptionSection}>
                            <h2 className={styles.sectionTitle}>Содержание</h2>
                            <MdRenderer source={project.content} />
                        </div>
                    )}

                    <div className={styles.slotsSection}>
                        <div className={styles.slotsSectionHeader}>
                            <h2 className={styles.sectionTitle}>Команда / Слоты</h2>
                            {isLeader && (
                                <Button
                                    variant="outline"
                                    onClick={() => { setIsCreatingSlot((previous) => !previous); }}
                                >
                                    <Plus size={16} />
                                    Добавить слот
                                </Button>
                            )}
                        </div>

                        {isCreatingSlot && (
                            <div className={styles.createSlotForm}>
                                <h3 className={styles.createSlotTitle}>Новый слот</h3>
                                <div className={styles.slotField}>
                                    <label className={styles.slotLabel} htmlFor="slot-title">Название позиции</label>
                                    <input
                                        id="slot-title"
                                        className={styles.slotInput}
                                        placeholder="Например: Backend Developer"
                                        value={newSlotTitle}
                                        onChange={(event) => { setNewSlotTitle(event.target.value); }}
                                        disabled={isSubmittingSlot}
                                    />
                                </div>
                                <div className={styles.slotField}>
                                    <SkillMultiSelect
                                        value={newSlotPrimarySkills}
                                        onChange={setNewSlotPrimarySkills}
                                        loadOptions={loadPrimarySkills}
                                        placeholder="Выберите основные навыки..."
                                        disabled={isSubmittingSlot}
                                        label="Основные навыки 1-го уровня"
                                    />
                                </div>
                                <div className={styles.slotField}>
                                    <SkillMultiSelect
                                        value={newSlotSecondarySkills}
                                        onChange={setNewSlotSecondarySkills}
                                        loadOptions={loadSecondarySkills}
                                        placeholder="Добавить доп. навыки..."
                                        disabled={isSubmittingSlot}
                                        max={10}
                                        label="Дополнительные навыки (2-й уровень, до 10)"
                                    />
                                </div>
                                <div className={styles.slotField}>
                                    <label className={styles.slotLabel} htmlFor="slot-desc">Описание</label>
                                    <textarea
                                        id="slot-desc"
                                        className={styles.slotTextarea}
                                        placeholder="Требования к кандидату..."
                                        value={newSlotDescription}
                                        onChange={(event) => { setNewSlotDescription(event.target.value); }}
                                        disabled={isSubmittingSlot}
                                    />
                                </div>
                                {slotCreateError !== undefined && (
                                    <p className={styles.slotError}>{slotCreateError}</p>
                                )}
                                <div className={styles.slotFormActions}>
                                    <Button
                                        variant="outline"
                                        onClick={() => { setIsCreatingSlot(false); }}
                                        disabled={isSubmittingSlot}
                                    >
                                        Отмена
                                    </Button>
                                    <Button
                                        onClick={() => { void handleCreateSlot(); }}
                                        disabled={isSubmittingSlot || !newSlotTitle || newSlotPrimarySkills.length === 0}
                                    >
                                        Создать слот
                                    </Button>
                                </div>
                            </div>
                        )}

                        {/* Skill filter dropdown for viewing slots by primary skill */}
                        {(project.slots ?? []).length > 0 && (
                            <div className={styles.skillFilterRow}>
                                <SkillSearch
                                    value={selectedSkillFilter}
                                    onChange={setSelectedSkillFilter}
                                    loadOptions={loadPrimarySkills}
                                    placeholder="Фильтр по основному навыку..."
                                />
                            </div>
                        )}

                        {(!project.slots || project.slots.length === 0) ? (
                            <p className={styles.emptySlots}>Слоты ещё не добавлены</p>
                        ) : (
                            <div className={styles.slotsList}>
                                {(project.slots
                                    .filter((slot: IProjectSlot) =>
                                        selectedSkillFilter === undefined
                                        || slot.primary_skills.some((s) => s.id === selectedSkillFilter.id)
                                    )
                                ).map((slot: IProjectSlot) => {
                                    const slotRequests = pendingRequestsBySlot(slot.id);
                                    const isOccupied = slot.user_id !== null;

                                    return (
                                        <div key={slot.id} className={styles.slotCard}>
                                            <div className={styles.slotHeader}>
                                                <div className={styles.slotInfo}>
                                                    <span className={styles.slotTitle}>{slot.title}</span>
                                                    {slot.primary_skills.map((skill) => (
                                                        <Badge key={skill.id} color={skill.color}>
                                                            {skill.name}
                                                        </Badge>
                                                    ))}
                                                    <Badge color={!isOccupied && slot.status === 'open' ? '039855' : undefined}>
                                                        {slotStatusLabel(isOccupied, slot.status)}
                                                    </Badge>
                                                </div>
                                                {isLeader && (
                                                    <button
                                                        type="button"
                                                        className={styles.deleteSlotBtn}
                                                        onClick={() => { setSlotIdToDelete(slot.id); }}
                                                        title="Удалить слот"
                                                    >
                                                        <X size={16} />
                                                    </button>
                                                )}
                                            </div>

                                            {/* Secondary skills badges */}
                                            {(slot.secondary_skills ?? []).length > 0 && (
                                                <div className={styles.secondarySkills}>
                                                    {(slot.secondary_skills ?? []).map((skill) => (
                                                        <Badge key={skill.id} color={skill.color}>
                                                            {skill.name}
                                                        </Badge>
                                                    ))}
                                                </div>
                                            )}

                                            {slot.description && (
                                                <p className={styles.slotDescription}>{slot.description}</p>
                                            )}

                                            {/* Apply button for authenticated non-leaders */}
                                            {userStore.isAuthenticated && !isLeader && !isOccupied && slot.status === 'open' && (
                                                <Button
                                                    variant="outline"
                                                    onClick={() => { void handleApply(slot.id); }}
                                                    disabled={applyingSlotId === slot.id}
                                                >
                                                    <UserCheck size={16} />
                                                    Откликнуться
                                                </Button>
                                            )}

                                            {/* Requests for leader */}
                                            {isLeader && slotRequests.length > 0 && (
                                                <div className={styles.requests}>
                                                    <p className={styles.requestsTitle}>Заявки ({slotRequests.length})</p>
                                                    {slotRequests.map((request) => (
                                                        <div key={request.id} className={styles.requestItem}>
                                                            <span className={styles.requestUser}>
                                                                Пользователь {request.user_id.slice(0, 8)}...
                                                            </span>
                                                            {request.cover_letter && (
                                                                <p className={styles.coverLetter}>{request.cover_letter}</p>
                                                            )}
                                                            <div className={styles.requestActions}>
                                                                <Button
                                                                    onClick={() => { void handleAcceptRequest(request.id); }}
                                                                >
                                                                    <Check size={16} />
                                                                    Принять
                                                                </Button>
                                                                <Button
                                                                    variant="outline"
                                                                    onClick={() => { void handleRejectRequest(request.id); }}
                                                                >
                                                                    <X size={16} />
                                                                    Отклонить
                                                                </Button>
                                                            </div>
                                                        </div>
                                                    ))}
                                                </div>
                                            )}
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </div>
                </div>

                <aside className={styles.sidebar}>
                    <div className={styles.leaderCard}>
                        <h3 className={styles.sidebarTitle}>Лидер проекта</h3>
                        {leader !== undefined && (
                            <div className={styles.leaderInfo}>
                                <UserAvatar
                                    avatarUrl={leader.avatar_url}
                                    nickname={leader.nickname}
                                    size={40}
                                />
                                <span className={styles.nickname}>{leader.nickname}</span>
                            </div>
                        )}
                    </div>

                    {project.idea_id !== null && (
                        <div className={styles.ideaCard}>
                            <h3 className={styles.sidebarTitle}>На основе идеи</h3>
                            <Link to="/idea/$ideaId" params={{ ideaId: project.idea_id }} className={styles.ideaLink}>
                                Перейти к идее
                            </Link>
                        </div>
                    )}
                </aside>
            </div>

            <ConfirmModal
                isOpen={isProjectDeleteModalOpen}
                title="Удалить проект?"
                description="Это действие необратимо. Все данные проекта, включая слоты и заявки, будут удалены навсегда."
                severity="danger"
                confirmLabel="Удалить"
                cancelLabel="Отмена"
                isLoading={isDeleting}
                onConfirm={() => { void handleDelete(); }}
                onCancel={() => { setIsProjectDeleteModalOpen(false); }}
            />
            <ConfirmModal
                isOpen={slotIdToDelete !== undefined}
                title="Удалить слот?"
                description="Слот и все связанные заявки будут удалены."
                severity="warning"
                confirmLabel="Удалить слот"
                cancelLabel="Отмена"
                onConfirm={() => {
                    if (slotIdToDelete !== undefined) {
                        void handleDeleteSlot(slotIdToDelete);
                    }
                }}
                onCancel={() => { setSlotIdToDelete(undefined); }}
            />
        </div>
    );
});
