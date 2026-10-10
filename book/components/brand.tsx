import Image from "next/image";

import logo from "../../assets/respkit.png";

export function Brand() {
  return <Image alt="RespKit" className="h-10 w-auto" priority src={logo} />;
}
