import { motion } from "framer-motion";
import { fadeUp } from "./animation";

type EmptyStateProps = {
  title: string;
  description: string;
};

export function EmptyState({ title, description }: EmptyStateProps) {
  return (
    <motion.div
      className="empty-state"
      initial="hidden"
      animate="visible"
      variants={fadeUp}
    >
      <h3>{title}</h3>
      <p>{description}</p>
    </motion.div>
  );
}
