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
import { CirclePlusIcon, Merge } from "lucide-react";
import React, { useEffect, useState } from "react";

const page = () => {
  const [rules, setRules] = useState([]);

  const sampleAST = {
    type: "operator",
    value: "AND",
    left: {
      type: "condition",
      value: ">",
      left: { type: "leaf", value: "age" },
      right: { type: "leaf", value: "30" },
    },
    right: {
      type: "condition",
      value: "=",
      left: { type: "leaf", value: "department" },
      right: { type: "leaf", value: "Sales" },
    },
  };

  let astData: any;

  // try {
  //   const ruleData = await apiHandler.get<any>(
  //     "/rules/671612f59437c7fa60122e2b"
  //   );

  //   astData = ruleData?.rule?.root_node;
  //   console.log(ruleData);

  //   console.log(astData);
  // } catch (error) {}

  useEffect(() => {
    const fetchRules = async () => {
      try {
        const rulesData = await apiHandler.get<any>("/rules");
        setRules(rulesData.rules);
      } catch (error) {
        console.error(error);
      }
    };

    fetchRules();
  }, []);

  return (
    <main className=" w-full max-w-[1350px] mx-auto my-5 flex flex-col justify-center items-center gap-5">
      <section className="w-full flex justify-between items-center ">
        <h1 className="text-5xl font-extrabold text- ">Dashboard</h1>
        <aside className="flex justify-center items-center gap-8">
          {/* <DrawerHandler isOpen={openDrawer}  title={"Create Rule"} />
           */}

          {/* <Button
            className=" flex justify-center items-center gap-2 px-4 py-2 "
            onClick={() => setopenDrawer(true)}
          >
            <CirclePlusIcon />
            Add Rule
          </Button>

          <DrawerHandler title="Are you sure?" isOpen={openDrawer} /> */}

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
      {rules && rules.length > 0 && (
        <div className="w-full flex flex-col justify-start items-start gap-4 ">
          <div className="text-3xl font-bold ">Rules :</div>
          {/* <section className="grid grid-cols-4 gap-4 w-full  rounded-md border border-gray-100/50 backdrop-blur-md  shadow-sm shadow-gray-100/25 p-4 "> */}
          <section className="grid grid-cols-3 gap-4 w-full p-4 ">
            {rules.map((rule, index) => (
              <RuleItem key={index} rule={rule} />
            ))}
          </section>
        </div>
      )}
      {/* <RuleForm /> */}

      {/* <ASTtree ast={astData} /> */}
    </main>
  );
};

export default page;
