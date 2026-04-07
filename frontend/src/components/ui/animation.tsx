import {
  AnimatePresence,
  motion,
  useReducedMotion,
  type HTMLMotionProps,
  type Transition,
  type Variants,
} from "framer-motion";
import type { PropsWithChildren, ReactNode } from "react";

const easeOut: [number, number, number, number] = [0.22, 1, 0.36, 1];

export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.12,
      delayChildren: 0.05,
    },
  },
};

export const fadeUp: Variants = {
  hidden: { opacity: 0, y: 24 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.62,
      ease: easeOut,
    },
  },
};

export const fadeLeft: Variants = {
  hidden: { opacity: 0, x: -28 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.62,
      ease: easeOut,
    },
  },
};

export const fadeRight: Variants = {
  hidden: { opacity: 0, x: 28 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.62,
      ease: easeOut,
    },
  },
};

export const scaleIn: Variants = {
  hidden: { opacity: 0, scale: 0.96, y: 12 },
  visible: {
    opacity: 1,
    scale: 1,
    y: 0,
    transition: {
      duration: 0.5,
      ease: easeOut,
    },
  },
};

export const routeTransition: Variants = {
  initial: { opacity: 0, y: 18, filter: "blur(6px)" },
  animate: {
    opacity: 1,
    y: 0,
    filter: "blur(0px)",
    transition: {
      duration: 0.42,
      ease: easeOut,
    },
  },
  exit: {
    opacity: 0,
    y: -12,
    filter: "blur(4px)",
    transition: {
      duration: 0.22,
      ease: "easeInOut",
    },
  },
};

export const cardHover = {
  y: -6,
  transition: {
    duration: 0.22,
    ease: easeOut,
  },
};

export const buttonHover = {
  y: -2,
  scale: 1.01,
  transition: {
    duration: 0.18,
    ease: easeOut,
  },
};

type RevealProps = PropsWithChildren<
  HTMLMotionProps<"div"> & {
    variants?: Variants;
    amount?: number;
    once?: boolean;
    onView?: boolean;
  }
>;

export function Reveal({
  children,
  variants = fadeUp,
  amount = 0.18,
  once = true,
  onView = true,
  ...props
}: RevealProps) {
  const reducedMotion = useReducedMotion();

  return (
    <motion.div
      initial="hidden"
      animate={onView ? undefined : "visible"}
      whileInView={onView ? "visible" : undefined}
      viewport={onView ? { once, amount } : undefined}
      variants={reducedMotion ? undefined : variants}
      {...props}
    >
      {children}
    </motion.div>
  );
}

type StaggerProps = PropsWithChildren<
  HTMLMotionProps<"div"> & {
    amount?: number;
    once?: boolean;
    onView?: boolean;
  }
>;

export function StaggerGroup({
  children,
  amount = 0.12,
  once = true,
  onView = true,
  ...props
}: StaggerProps) {
  const reducedMotion = useReducedMotion();

  return (
    <motion.div
      initial="hidden"
      animate={onView ? undefined : "visible"}
      whileInView={onView ? "visible" : undefined}
      viewport={onView ? { once, amount } : undefined}
      variants={reducedMotion ? undefined : staggerContainer}
      {...props}
    >
      {children}
    </motion.div>
  );
}

type RouteTransitionProps = {
  children: ReactNode;
  routeKey: string;
};

export function RouteTransition({ children, routeKey }: RouteTransitionProps) {
  const reducedMotion = useReducedMotion();

  if (reducedMotion) {
    return <>{children}</>;
  }

  return (
    <AnimatePresence mode="wait" initial={false}>
      <motion.div
        key={routeKey}
        variants={routeTransition}
        initial="initial"
        animate="animate"
        exit="exit"
      >
        {children}
      </motion.div>
    </AnimatePresence>
  );
}

type PresenceMessageProps = {
  children: ReactNode;
  className?: string;
};

export function PresenceMessage({ children, className }: PresenceMessageProps) {
  return (
    <AnimatePresence mode="popLayout">
      {children ? (
        <motion.div
          key="message"
          className={className}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -6 }}
          transition={{ duration: 0.22, ease: easeOut }}
        >
          {children}
        </motion.div>
      ) : null}
    </AnimatePresence>
  );
}

export function delayedTransition(delay = 0): Transition {
  return {
    duration: 0.62,
    delay,
    ease: easeOut,
  };
}
