"use client ";

import { ArrowDown, ArrowUp, Cloud, Thermometer } from "lucide-react";
import React, { useEffect, useState } from "react";
import SummaryCard from "./SummaryCard";
import apiHandler from "@/handlers/apiHandler";
import { WeatherSummary } from "@/lib/types";

type Props = {
  position: string;
};

const SummaryCardInfo: {
  title: string;
  icon: React.ReactElement;
  color: string;
  stat: keyof WeatherSummary;
}[] = [
  {
    title: "Average Temperature",
    icon: <Thermometer className=" text-blue-500" size={30} />,
    color: "#3B82F6",
    stat: "avg_temperature",
  },
  {
    title: "Maximum Temperature",
    icon: <ArrowUp className=" text-red-500" size={30} />,
    color: "#EF4444",
    stat: "max_temperature",
  },
  {
    title: "Minimum Temperature",
    icon: <ArrowDown className=" text-green-500" size={30} />,
    color: "#10B981",
    stat: "min_temperature",
  },
  {
    title: "Dominant Condition ",
    icon: <Cloud className=" text-indigo-500" size={30} />,
    color: "#6366f1",
    stat: "dominant_condition",
  },
];

const DailySummary = ({ position }: Props) => {
  const [summaryData, setsummaryData] = useState<WeatherSummary>();

  useEffect(() => {
    const DailySummaryData = async () => {
      try {
        console.log("get summary api is calling for", position);
        const summary = await apiHandler.get<WeatherSummary>(
          `getsummary/${"Mumbai"}`
        );

        console.log("summaryData", summary);
        setsummaryData(summary);
      } catch (error) {
        console.log("error", error);
      }
    };
    DailySummaryData();
  }, [position]);

  return (
    <section className="w-full flex flex-col justify-start items-start  gap-3">
      <div className="w-full text-4xl font-bold">Daily Summary</div>

      <div className="w-full flex justify-center items-center gap-4">
        {summaryData ? (
          SummaryCardInfo.map((card) => {
            return (
              <div>
                <SummaryCard summaryData={summaryData} card={card} />
              </div>
            );
          })
        ) : (
          <div>Wait data is loading...</div>
        )}
      </div>
    </section>
  );
};

export default DailySummary;
