import type { JSX } from "react";

import { useTheme } from "@/shared/lib/hooks/useTheme";

export interface LogoProps {
  className?: string | undefined;
  height?: number | undefined;
}

export function Logo({ className, height = 36 }: LogoProps): JSX.Element {
  const { theme } = useTheme();
  const source = theme === "dark" ? "/devspace-logo-dark.png" : "/DevSpaceLogo-removebg.png";

  return <img src={source} alt="DevSpace" className={className} height={height} />;
}
