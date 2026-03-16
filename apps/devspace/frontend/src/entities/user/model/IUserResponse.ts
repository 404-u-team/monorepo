export interface IUserResponse {
    id: string;
    nickname: string;
    avatar_uri: string;
    bio: string;
    main_role: string;
    skills: ISkill[];
}

export interface ISkill {
    id: string;
    name: string;
    parent_id: string;
    children: string[];
}
