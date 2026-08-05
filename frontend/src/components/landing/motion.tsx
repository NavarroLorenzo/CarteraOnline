import { motion, type HTMLMotionProps, type Variants } from "framer-motion";
import type { PropsWithChildren } from "react";
import {
  buttonHover,
  cardHover,
  delayedTransition,
  fadeLeft,
  fadeRight,
  fadeUp,
  scaleIn,
  staggerContainer,
} from "../ui/animation";

export {
  buttonHover,
  cardHover,
  delayedTransition,
  fadeLeft,
  fadeRight,
  fadeUp,
  scaleIn,
  staggerContainer,
};

type RevealProps = PropsWithChildren<
  HTMLMotionProps<"div"> & {
    variants?: Variants;
    amount?: number;
    once?: boolean;
    onView?: boolean;
  }
>;

/**
 * Landing-specific reveal.
 *
 * The homepage preview is an animated product demonstration, so it keeps its
 * own predictable motion behavior. Application screens continue using the
 * accessibility-aware shared wrapper.
 */
export function Reveal({
  children,
  variants = fadeUp,
  amount = 0.18,
  once = true,
  onView = true,
  ...props
}: RevealProps) {
  return (
    <motion.div
      initial="hidden"
      animate={onView ? undefined : "visible"}
      whileInView={onView ? "visible" : undefined}
      viewport={onView ? { once, amount } : undefined}
      variants={variants}
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
  return (
    <motion.div
      initial="hidden"
      animate={onView ? undefined : "visible"}
      whileInView={onView ? "visible" : undefined}
      viewport={onView ? { once, amount } : undefined}
      variants={staggerContainer}
      {...props}
    >
      {children}
    </motion.div>
  );
}
