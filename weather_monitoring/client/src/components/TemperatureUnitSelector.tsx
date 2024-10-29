"use client";
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
import { Units } from "@/lib/types";

type Props = {
  onUnitChange: (position: Units) => void;
};

const TemperatureUnitSelector = ({ onUnitChange }: Props) => {
  const [unit, setUnit] = React.useState<Units>("celsius"); // Initial positionas");

  const handlePositionChange = (newPosition: Units) => {
    setUnit(newPosition);
    onUnitChange(newPosition); // Call the parent callback with the new position
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline">Temperature Unit </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56">
        <DropdownMenuLabel>Select Unit</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuRadioGroup
          value={unit}
          onValueChange={handlePositionChange}
        >
          <DropdownMenuRadioItem value="celsius">Celcius</DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="fahrenheit">
            Fahrenheit
          </DropdownMenuRadioItem>
          <DropdownMenuRadioItem value="kelvin">Kelvin</DropdownMenuRadioItem>
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default TemperatureUnitSelector;
