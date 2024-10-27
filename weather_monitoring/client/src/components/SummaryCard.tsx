import React, { ReactElement } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { WeatherSummary } from "@/lib/types";

type Props = {
  card: {
    title: string;
    color: string;
    icon: ReactElement;
  };
  summaryData: any;
};

function SummaryCard({ card, summaryData }: Props) {
  //@ts-ignore
  const statValue = summaryData[card?.stat as keyof WeatherSummary];
  const roundedStatValue =
    typeof statValue === "number" ? statValue.toFixed(2) : statValue;
  return (
    <Card className="bg-gray-100/75 w-[350px] h-[250px]">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className=" text-gray-900">{card.title}</CardTitle>
        <span>{card.icon}</span>
      </CardHeader>
      <CardContent className=" w-full  flex  justify-center items-center">
        <div className="text-4xl font-bold" style={{color : card.color}}>
          {roundedStatValue}
        </div>
      </CardContent>
    </Card>
  );
}

export default SummaryCard;
