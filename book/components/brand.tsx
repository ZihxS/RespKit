import Image from "next/image";

import logo from "../../assets/respkit.png";
import { brandName } from "@/lib/branding";

export function Brand() {
  return <Image alt={brandName} className="h-10 w-auto" priority src={logo} />;
}
