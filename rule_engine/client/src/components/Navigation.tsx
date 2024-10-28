import React from "react";
import ModeToggle from "./ModeToggle";
import Image from "next/image";
import logo from "../public/logo.svg";

type Props = {};

const Navigation = (props: Props) => {
  return (
    <div className="p-4 flex items-center justify-between relative mx-2">
      <aside className="flex items-center gap-2 ">
        <Image src="/logo.svg" alt="" width={40} height={40} />
        <h4 className="sm:text-4xl font-bold text-primary">Rule Engine</h4>
      </aside>

      <aside className="flex gap-2 items-center">
        <ModeToggle />
      </aside>
    </div>
  );
};

export default Navigation;
