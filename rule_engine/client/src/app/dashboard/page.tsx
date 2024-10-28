"use client";
import ASTtree from "@/components/AST-tree";
import DrawerHandler from "@/components/Drawer";
import CombineRuleForm from "@/components/forms/Combine-rule-form";
import RuleForm from "@/components/forms/Rule-form";
import Model from "@/components/Model";
import RuleItem from "@/components/RuleItem";
import Tooltiper from "@/components/Tooltip";
import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { Separator } from "@/components/ui/separator";
import apiHandler from "@/handlers/api-handler";
import { CirclePlusIcon, Loader2, Merge } from "lucide-react";
import React, { useEffect, useState } from "react";

const page = () => {
  const [rules, setRules] = useState([]);
  const [isLoading, setisLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchRules = async () => {
      try {
        setisLoading(true);
        const rulesData = await apiHandler.get<any>("/rules");
        console.log("rulesData", rulesData);
        setRules(rulesData.rules);
        setisLoading(false);
      } catch (error) {
        console.error(error);
        setisLoading(false);
      }
    };

    fetchRules();
  }, []);

  return (
    <main className=" w-full max-w-[1350px] mx-auto my-5 flex flex-col justify-center items-center gap-5">
      <section className="w-full flex justify-between items-center ">
        <h1 className="text-5xl font-extrabold text- ">Dashboard</h1>
        <aside className="flex justify-center items-center gap-8">
          <Drawer modal={false}>
            <DrawerTrigger asChild>
              <Button className=" flex justify-center items-center gap-2 px-4 py-2 ">
                <CirclePlusIcon />
                Add Rule
              </Button>
            </DrawerTrigger>
            <DrawerContent className="w-full flex flex-col  justify-center items-center ">
              <DrawerHeader>
                <DrawerTitle>Create a Rule</DrawerTitle>
                <DrawerDescription>
                  Please fill the Required information
                </DrawerDescription>
              </DrawerHeader>
              <RuleForm />
              <DrawerFooter>
                <DrawerClose asChild>
                  <Button variant="outline">Close</Button>
                </DrawerClose>
              </DrawerFooter>
            </DrawerContent>
          </Drawer>
          <Model
            desc={`You can combine more than 2 rules to create a more complex rule.  `}
            title="Combine rules"
            TriggerElement={
              <Tooltiper tooltipMessage="Add more than 2 rules ">
                <span className="flex justify-center items-center  gap-4 ">
                  <Merge />
                  <div>Combine Rule</div>
                </span>
              </Tooltiper>
            }
            isforAST={false}
          >
            <CombineRuleForm />
          </Model>
        </aside>
      </section>
      <Separator className="w-full" />
      {!isLoading ? (
        rules && rules.length > 0 ? (
          <div className="w-full flex flex-col justify-start items-start gap-4 ">
            <div className="text-3xl font-bold ">Rules :</div>
            <section className="grid grid-cols-3 gap-4 w-full p-4 ">
              {rules.map((rule, index) => (
                <RuleItem key={index} rule={rule} />
              ))}
            </section>
          </div>
        ) : (
          <div className="w-full flex justify-center items-center">
            <h2 className="text-4xl font-bold">No Rules data available</h2>
          </div>
        )
      ) : (
        <div className="w-full  flex justify-center items-center animate-spin ">
          <Loader2 size={100} />
        </div>
      )}
    </main>
  );
};

export default page;
