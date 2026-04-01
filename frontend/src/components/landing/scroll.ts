const easeOutExpo = (progress: number): number => {
  if (progress >= 1) {
    return 1;
  }

  return 1 - Math.pow(2, -10 * progress);
};

export function scrollToSection(sectionId: string, offset = 104): void {
  if (typeof window === "undefined" || typeof document === "undefined") {
    return;
  }

  const section = document.getElementById(sectionId);

  if (!section) {
    return;
  }

  const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  const startY = window.scrollY;
  const targetY = Math.max(0, section.getBoundingClientRect().top + startY - offset);

  if (prefersReducedMotion) {
    window.scrollTo(0, targetY);
    return;
  }

  const distance = targetY - startY;
  const duration = 820;
  let animationStart: number | null = null;

  const step = (timestamp: number) => {
    if (animationStart === null) {
      animationStart = timestamp;
    }

    const elapsed = timestamp - animationStart;
    const progress = Math.min(elapsed / duration, 1);
    const easedProgress = easeOutExpo(progress);

    window.scrollTo(0, startY + distance * easedProgress);

    if (progress < 1) {
      window.requestAnimationFrame(step);
    }
  };

  window.requestAnimationFrame(step);
}
