"use client";
import React, { createContext, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import TempThresholdForm from "./form/Temp-threshold-form";
import { Button } from "./ui/button";

type Props = {};

interface DialogContextType {
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}

export const DialogContext = createContext<DialogContextType>({
  setOpen: () => {},
});

const Tempthreshold = (props: Props) => {
  const [open, setOpen] = useState<boolean>(false);
  return (
    <DialogContext.Provider value={{ setOpen }}>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger className="w-[100px] border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground py-1 rounded-md ">
          <span className="text-sm ">Set threshold</span>
        </DialogTrigger>
        <DialogContent className="w-[1000px]">
          <DialogHeader className="w-full">
            <DialogTitle>Set Threshold</DialogTitle>
            <DialogDescription>
              You can set the temperature threshold for each city.You will get
              the email notifications if temperature consecutively cross your
              set temperature threshold.
            </DialogDescription>
          </DialogHeader>
          <TempThresholdForm />
        </DialogContent>
      </Dialog>
    </DialogContext.Provider>
  );
};

export default Tempthreshold;
