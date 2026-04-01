import {
  motion,
  useReducedMotion,
  type HTMLMotionProps,
  type Transition,
  type Variants,
} from "framer-motion";
import type { PropsWithChildren } from "react";

const easeOut: [number, number, number, number] = [0.22, 1, 0.36, 1];

export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.14,
      delayChildren: 0.08,
    },
  },
};

export const fadeUp: Variants = {
  hidden: { opacity: 0, y: 28 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.72,
      ease: easeOut,
    },
  },
};

export const fadeLeft: Variants = {
  hidden: { opacity: 0, x: -32 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.72,
      ease: easeOut,
    },
  },
};

export const fadeRight: Variants = {
  hidden: { opacity: 0, x: 32 },
  visible: {
    opacity: 1,
    x: 0,
    transition: {
      duration: 0.72,
      ease: easeOut,
    },
  },
};

export const cardHover = {
  y: -8,
  transition: {
    duration: 0.22,
    ease: easeOut,
  },
};

export const buttonHover = {
  y: -2,
  scale: 1.01,
  transition: {
    duration: 0.2,
    ease: easeOut,
  },
};

type RevealProps = PropsWithChildren<
  HTMLMotionProps<"div"> & {
    variants?: Variants;
    amount?: number;
    once?: boolean;
  }
>;

export function Reveal({
  children,
  variants = fadeUp,
  amount = 0.2,
  once = true,
  ...props
}: RevealProps) {
  const reducedMotion = useReducedMotion();

  return (
    <motion.div
      initial="hidden"
      whileInView="visible"
      viewport={{ once, amount }}
      variants={reducedMotion ? undefined : variants}
      {...props}
    >
      {children}
    </motion.div>
  );
}

type StaggerProps = PropsWithChildren<HTMLMotionProps<"div">>;

export function StaggerGroup({ children, ...props }: StaggerProps) {
  const reducedMotion = useReducedMotion();

  return (
    <motion.div
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, amount: 0.16 }}
      variants={reducedMotion ? undefined : staggerContainer}
      {...props}
    >
      {children}
    </motion.div>
  );
}

export function delayedTransition(delay = 0): Transition {
  return {
    duration: 0.72,
    delay,
    ease: easeOut,
  };
}
