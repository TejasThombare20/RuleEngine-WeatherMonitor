"use client";
import CitySelector from "@/components/City-Selector";
import DailySummary from "@/components/DailySummary";
import Tempthreshold from "@/components/Temp-threshold";
import TemperatureLineChart from "@/components/Temperature-Line-chart";
import { PositionType } from "@/lib/types";
import React, { useState } from "react";

type Props = {};

const page = (props: Props) => {
  const [position, onPositionChange] = useState<PositionType>("allcities");
  return (
    <main className="max-w-[1450px] mx-auto my-2 flex flex-col justify-center items-center overflow-hidden gap-4">
      <section className="w-full flex justify-between items-center ">
        <div className="w-full flex flex-col justify-start items-start gap-4">
          <h2 className="text-5xl font-bold">Dashboard</h2>

        </div>
        <aside className="flex justify-center items-center gap-4">
           <Tempthreshold/>
          <CitySelector onPositionChange={onPositionChange} />
        </aside>
      </section>
      <section className="w-full max-w-[1400px] mx-auto">
        <TemperatureLineChart position={position} />
      </section>
      {position != "allcities" ? (
        <section>
          <DailySummary position={position} />
        </section>
      ) : (
        <section className="w-full text-3xl flex justify-center items-center font-bold">
          Please select a city from top right corner to view its daily summary.
        </section>
      )}
    </main>
  );
};

export default page;
