"use client";
import React, { useEffect, useState } from "react";
import { LineChart } from "./Line-chart";
import { PositionType, Units } from "@/lib/types";
import apiHandler from "@/handlers/apiHandler";
import {
  addHours,
  endOfDay,
  format,
  getHours,
  parseISO,
  startOfDay,
} from "date-fns";
import { Loader2 } from "lucide-react";

type Props = {
  isSingleCity?: boolean;
  position: PositionType;
  unit: Units;
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

const getUnitSymbol = (unit: Units): string => {
  switch (unit) {
    case "fahrenheit":
      return "°F";
    case "kelvin":
      return "K";
    case "celsius":
    default:
      return "°C";
  }
};

const convertTemperature = (celsius: number, targetUnit: Units): number => {
  switch (targetUnit) {
    case "fahrenheit":
      return (celsius * 9) / 5 + 32;
    case "kelvin":
      return celsius + 273.15;
    case "celsius":
    default:
      return celsius;
  }
};

const processTimeSeriesData = (data: TemperatureData[], unit: Units) => {
  const currentDate = new Date();
  const dayStart = startOfDay(currentDate);
  const dayEnd = endOfDay(currentDate);

  const latestTimeStamp = data?.[0]?.timestamp;
  const lastesthour = latestTimeStamp
    ? getHours(new Date(latestTimeStamp))
    : null;

  // Create array of 4-hour interval timestamps
  const timeLabels = [];
  for (let hour = 0; hour <= 24; hour += 4) {
    if (hour > lastesthour!) {
      break;
    }
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
        // Convert temperature to the selected unit
        dataPoint[city] = convertTemperature(closestReading.temperature, unit);
      }
    });

    return dataPoint;
  });

  return formattedData;
};

const TemperatureLineChart = ({
  position,
  isSingleCity = false,
  unit,
}: Props) => {
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
    if (data && data?.length > 0) {
      const procData = processTimeSeriesData(data, unit);
      setProcessedData(procData);
    }
  }, [data, unit]);

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
          className=" w-full h-96"
          data={processedData}
          index="time"
          showXAxis={true}
          showYAxis={true}
          enableLegendSlider={true}
          valueFormatter={(value: number) => `${value.toFixed(1)}${getUnitSymbol(unit)}`}
          categories={categories}
          tickGap={1}
          xAxisLabel="X-axis"
        />
      )}
    </div>
  );
};

export default TemperatureLineChart;
