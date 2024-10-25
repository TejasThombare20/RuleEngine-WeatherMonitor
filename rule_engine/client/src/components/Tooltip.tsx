import React from "react";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "./ui/tooltip";

type Props = {
  children: React.ReactElement;
  tooltipMessage: string;
};

const Tooltiper = ({ children, tooltipMessage }: Props) => {
  return (
    <TooltipProvider>
      <Tooltip >
        <TooltipTrigger asChild>{children}</TooltipTrigger>
        <TooltipContent className="bg-background dark:text-gray-100 text-gray-800">
          <p>{tooltipMessage}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default Tooltiper;
