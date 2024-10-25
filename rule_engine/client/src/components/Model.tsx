import React from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Button } from "./ui/button";
import { Network } from "lucide-react";
import clsx from "clsx";

type Props = {
  children: React.ReactNode;
  title: string;
  desc: string;
  TriggerElement: React.ReactElement;
  isforAST: boolean;
};

const Model = ({ children, desc, title, TriggerElement, isforAST }: Props) => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">{TriggerElement}</Button>
      </DialogTrigger>
      <DialogContent
        className={clsx(
          isforAST && "max-w-[1400px] h-[750px]",
          !isforAST && " w-[400px]"
        )}
      >
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            <div>{desc}</div>
          </DialogDescription>
        </DialogHeader>
        {children}
        <DialogFooter className="sm:justify-start">
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              Close
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default Model;
