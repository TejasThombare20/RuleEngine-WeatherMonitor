"use client";
import React, { useEffect, useState } from "react";
import { LineChart } from "./Line-chart";
import { PositionType } from "@/lib/types";
import apiHandler from "@/handlers/apiHandler";
import { addHours, endOfDay, format, parseISO, startOfDay } from "date-fns";
import { Loader2 } from "lucide-react";

type Props = {
  isSingleCity?: boolean;
  position: PositionType;
};

interface TemperatureData {
  city_name: string;
  temperature: number;
  feels_like: number;
  condition: string;
  timestamp: string;
}

interface FormattedData {
  time: string;
  [key: string]: number | string;
}

const processTimeSeriesData = (data: TemperatureData[]) => {
  const currentDate = new Date();
  const dayStart = startOfDay(currentDate);
  const dayEnd = endOfDay(currentDate);

    
    // Create array of 4-hour interval timestamps
    const timeLabels = [];
    for (let hour = 0; hour <= 24; hour += 4) {
        timeLabels.push(format(addHours(dayStart, hour), "HH:mm"));
    }

  // Format data for chart
  const formattedData = timeLabels.map((timeLabel) => {
    const dataPoint: FormattedData = { time: timeLabel };

    // For each city, find the closest temperature reading to this time
    const uniqueCities = Array.from(new Set(data.map((d) => d.city_name)));

    uniqueCities.forEach((city) => {
      const cityData = data.filter((d) => d.city_name === city);
      if (cityData.length > 0) {
        // Find closest reading
        const timeHour = parseInt(timeLabel.split(":")[0]);
        const closestReading = cityData.reduce((prev, curr) => {
          const currHour = new Date(curr.timestamp).getHours();
          const prevHour = new Date(prev.timestamp).getHours();
          return Math.abs(currHour - timeHour) < Math.abs(prevHour - timeHour)
            ? curr
            : prev;
        });
        dataPoint[city] = closestReading.temperature;
      }
    });

    return dataPoint;
  });

  return formattedData;
};

const TemperatureLineChart = ({ position, isSingleCity = false }: Props) => {
  const [data, setdata] = useState<TemperatureData[]>([]);
  const [isLaoding, setisLaoding] = useState(false);
  const [processedData, setProcessedData] = useState<FormattedData[]>([]);
  useEffect(() => {
    const fetchCityData = async () => {
      setisLaoding(true);
      try {
        const weatherData = await apiHandler.get<TemperatureData[]>(
          `/getRecords/${position}`
        );
        console.log("weatherData", weatherData);

        setdata(weatherData);
        setisLaoding(false);
      } catch (error) {
        setisLaoding(false);
        console.log(error);
      }
    };
    fetchCityData();
  }, [position]); // only run once when component mounts

  useEffect(() => {
    if ( data && data?.length > 0) {
      const procData = processTimeSeriesData(data);
      setProcessedData(procData);
    }
  }, [data]);

  // Get categories based on whether it's single city or multi-city view
  const categories = isSingleCity
    ? ["temperature"]
    : Array.from(new Set(data?.map((entry) => entry.city_name)));

  return (
    <div className="w-full my-5 h-[400px]">
      {isLaoding ? (
        <div className="w-full max-h-screen flex justify-center items-center">
          <Loader2 className="animate-spin" size={200} />
        </div>
      ) : (
        <LineChart
          className=" w-full h-80"
          data={processedData}
          index="time"
          showXAxis={true}
          showYAxis={true}
          enableLegendSlider={true}
          valueFormatter={(value: number) => `${value.toFixed(1)}°C`}
          categories={categories}
          tickGap={1}
          xAxisLabel="X-axis"
        />
      )}
    </div>
  );
};

export default TemperatureLineChart;
