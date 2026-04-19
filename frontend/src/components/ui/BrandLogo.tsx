import logoFull from "../../assets/cenz_logo_entero.png";
import logoIcon from "../../assets/cenz_logo_solo.png";

type BrandLogoVariant = "full" | "icon" | "responsive";

type BrandLogoProps = {
  variant?: BrandLogoVariant;
  className?: string;
};

export function BrandLogo({ variant = "full", className }: BrandLogoProps) {
  const baseClassName = `brand-logo brand-logo--${variant}${className ? ` ${className}` : ""}`;

  if (variant === "responsive") {
    return (
      <span className={baseClassName} role="img" aria-label="Cenz">
        <img className="brand-logo__image brand-logo__image--full" src={logoFull} alt="" aria-hidden="true" />
        <img className="brand-logo__image brand-logo__image--icon" src={logoIcon} alt="" aria-hidden="true" />
      </span>
    );
  }

  return (
    <span className={baseClassName}>
      <img
        className="brand-logo__image"
        src={variant === "icon" ? logoIcon : logoFull}
        alt="Cenz"
      />
    </span>
  );
}
