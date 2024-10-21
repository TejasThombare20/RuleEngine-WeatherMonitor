import ASTtree from "@/components/AST-tree";
import RuleForm from "@/components/forms/rule-form";
import apiHandler from "@/handlers/api-handler";
import React, { useEffect } from "react";

const page = async () => {
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

  try {
    const ruleData = await apiHandler.get<any>(
      "/rules/671612f59437c7fa60122e2b"
    );

    astData = ruleData?.rule?.root_node;
    console.log(ruleData);

    console.log(astData);
  } catch (error) {}

  return (
    <div className="w-full h-full flex  justify-center items-center">
      {/* <RuleForm /> */}

      <ASTtree ast={astData} />
    </div>
  );
};

export default page;
