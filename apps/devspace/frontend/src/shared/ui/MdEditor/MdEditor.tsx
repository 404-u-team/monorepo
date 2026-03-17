import type { JSX } from 'react';
import { clsx } from 'clsx';
import MDEditor, { commands } from '@uiw/react-md-editor';
import '@uiw/react-md-editor/markdown-editor.css';
import '@uiw/react-markdown-preview/markdown.css';
import styles from './MdEditor.module.scss';

const titleGroupCommand = commands.group(
    [
        commands.heading1,
        commands.heading2,
        commands.heading3,
        commands.heading4,
        commands.heading5,
        commands.heading6,
    ],
    {
        name: 'title',
        groupName: 'title',
        buttonProps: { 'aria-label': 'Заголовки' },
    },
);

const editorCommands = [
    commands.bold,
    commands.italic,
    commands.strikethrough,
    commands.divider,
    titleGroupCommand,
    commands.divider,
    commands.quote,
    commands.code,
    commands.codeBlock,
    commands.divider,
    commands.unorderedListCommand,
    commands.orderedListCommand,
    commands.checkedListCommand,
    commands.divider,
    commands.link,
    commands.image,
];

const extraEditorCommands = [
    commands.codeEdit,
    commands.codeLive,
    commands.codePreview,
    commands.divider,
    commands.fullscreen,
];

export interface MdEditorProps {
    id?: string | undefined;
    value: string;
    onChange: (value: string) => void;
    height?: number | undefined;
    minHeight?: number | undefined;
    placeholder?: string | undefined;
    disabled?: boolean | undefined;
    className?: string | undefined;
}

export function MdEditor({
    id,
    value,
    onChange,
    height = 360,
    minHeight,
    placeholder,
    disabled = false,
    className,
}: MdEditorProps): JSX.Element {
    return (
        <div className={clsx(styles.wrapper, className)} data-color-mode="light">
            <MDEditor
                value={value}
                onChange={(v) => { if (!disabled) { onChange(v ?? ''); } }}
                commands={editorCommands}
                extraCommands={extraEditorCommands}
                height={height}
                {...(minHeight !== undefined ? { minHeight } : {})}
                visibleDragbar={false}
                textareaProps={{
                    id,
                    placeholder,
                    disabled,
                }}
            />
        </div>
    );
}
