import MDEditor, { commands, type ICommand } from "@uiw/react-md-editor";
import { clsx } from "clsx";
import type { JSX } from "react";

import "@uiw/react-md-editor/markdown-editor.css";
import "@uiw/react-markdown-preview/markdown.css";
import { useTheme } from "@/shared/lib/hooks/useTheme";

import styles from "./MdEditor.module.scss";

function withTitle(command: ICommand, title: string): ICommand {
  return { ...command, buttonProps: { ...command.buttonProps, title } };
}

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
    name: "title",
    groupName: "title",
    buttonProps: { "aria-label": "Заголовки", title: "Заголовки" },
  },
);

const editorCommands = [
  withTitle(commands.bold, "Жирный (Ctrl+B)"),
  withTitle(commands.italic, "Курсив (Ctrl+I)"),
  withTitle(commands.strikethrough, "Зачёркнутый"),
  commands.divider,
  titleGroupCommand,
  commands.divider,
  withTitle(commands.quote, "Цитата"),
  withTitle(commands.code, "Код"),
  withTitle(commands.codeBlock, "Блок кода"),
  commands.divider,
  withTitle(commands.unorderedListCommand, "Ненумерованный список"),
  withTitle(commands.orderedListCommand, "Нумерованный список"),
  withTitle(commands.checkedListCommand, "Список задач"),
  commands.divider,
  withTitle(commands.link, "Ссылка"),
  withTitle(commands.image, "Изображение"),
];

const extraEditorCommands = [
  withTitle(commands.codeEdit, "Редактор"),
  withTitle(commands.codeLive, "Просмотр вживую"),
  withTitle(commands.codePreview, "Предпросмотр"),
  commands.divider,
  withTitle(commands.fullscreen, "Полноэкранный режим"),
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
  const { theme } = useTheme();
  return (
    <div className={clsx(styles.wrapper, className)} data-color-mode={theme}>
      <MDEditor
        value={value}
        onChange={(v) => {
          if (!disabled) {
            onChange(v ?? "");
          }
        }}
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
