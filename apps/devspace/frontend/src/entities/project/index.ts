export { ProjectCard, type ProjectCardProps } from './ui/ProjectCard/ProjectCard'
export { ProjectCardSkeleton } from './ui/ProjectCardSkeleton/ProjectCardSkeleton'
export type { IProject } from './model/IProject'
export type { IProjectSlot, IProjectSlotSkill } from './model/IProjectSlot'
export type { IRequest } from './model/IRequest'
export {
    fetchProjectById,
    fetchProjectSlots,
    fetchProjects,
    createProject,
    updateProject,
    deleteProject,
    createProjectSlot,
    deleteProjectSlot,
    updateProjectSlot,
    applyToSlot,
    getProjectRequests,
    acceptRequest,
    rejectRequest,
    getMyRequests,
} from './api/projectApi'
export type { IProjectDetailResponse, FetchProjectsParameters, PaginatedProjects } from './api/projectApi'
