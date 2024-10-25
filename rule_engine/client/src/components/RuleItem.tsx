import React from "react";
import { Card, CardContent } from "./ui/card";
import { Network, ScanLine } from "lucide-react";
import Model from "./Model";
import ASTtree from "./AST-tree";
import DrawerHandler from "./Drawer";
import Verifyruleform from "./forms/Verify-rule-form";

type Props = {
  rule: any;
};

const RuleItem = ({ rule }: Props) => {
  return (
    <Card className="relative py-2 w-[400px] h-[250px] rounded-md border border-gray-100/50 backdrop-blur-sm shadow-sm shadow-gray-100/25 ">
      <CardContent className="flex flex-col justify-start items-start h-full">
        <p className="text-base">
          <span>Rule name: </span>
          {rule.name}
        </p>
        <p>Description:</p>
        <pre className=" w-full h-full bg-gray-100/60 text-gray-950 p-1 rounded-md overflow-y-scroll ">
          <p className="tracking-wider  text-wrap ">{rule.description}</p>
        </pre>
      </CardContent>
      <span className=" absolute flex justify-center items-center top-2 right-1">
        <Model
          TriggerElement={<Network />}
          title="AST Tree"
          desc={`AST Tree for rule :  ${rule.name} `}
          isforAST={true}
        >
          <ASTtree ast={rule?.root_node} />
        </Model>
        <Model
          desc={`Evaulate your JSON data against rule :  ${rule.name} `}
          title="Rule Verification"
          TriggerElement={<ScanLine />}
          isforAST={false}
        >
          <Verifyruleform rule_id={rule?._id} />
        </Model>
      </span>
    </Card>
  );
};

export default RuleItem;
