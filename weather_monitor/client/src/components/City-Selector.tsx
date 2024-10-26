import React from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Button } from "./ui/button";
import { PositionType } from "@/lib/types";

type Props = {
  onPositionChange: (position: PositionType) => void;
};

const CitySelector = ({ onPositionChange }: Props) => {
  const [position, setPosition] = React.useState<PositionType>("allcities"); // Initial positionas");

  const handlePositionChange = (newPosition: PositionType) => {
    setPosition(newPosition);
    onPositionChange(newPosition); // Call the parent callback with the new position
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline">Select City </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        <DropdownMenuLabel>Select City</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuRadioGroup
          value={position}
          onValueChange={handlePositionChange}
        >
          <DropdownMenuRadioItem value="allcities">
            All Cities
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="Mumbai">Mumbai</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="Delhi">Delhi</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="Bengluru">
            Bengluru
          </DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default CitySelector;
