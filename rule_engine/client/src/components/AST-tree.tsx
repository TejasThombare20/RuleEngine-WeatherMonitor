"use client";

import React, { useCallback, useEffect } from "react";

import {
  addEdge,
  Connection,
  Edge,
  NodeTypes,
  ReactFlow,
  useEdgesState,
  useNodesState,
  Node,
  EdgeProps,
  getBezierPath,
  Controls,
  MiniMap,
  getStraightPath,
  BaseEdge,
  NodeProps,
  Handle,
  Position,
} from "@xyflow/react";

// you also need to adjust the style import
import "@xyflow/react/dist/style.css";

// or if you just want basic styles
import "@xyflow/react/dist/base.css";

interface ASTNode {
  type: "operator" | "condition" | "leaf" | String;
  value: string;
  left?: ASTNode;
  right?: ASTNode;
}

const CustomNode: React.FC<NodePropss> = ({ data }) => {
  return (
    <div className="px-4 py-2 shadow-md rounded-md bg-white border-2 border-gray-200">
      <Handle
        type="target"
        position={Position.Top}
        id="top"
        style={{ background: "#555" }}
      />
      <div className="font-bold">{data.label}</div>
      <Handle
        type="source"
        position={Position.Bottom}
        id="bottom"
        style={{ background: "#555" }}
      />
    </div>
  );
};

// function CustomEdge({ id, sourceX, sourceY, targetX, targetY }) {
//   const [edgePath] = getStraightPath({
//     sourceX,
//     sourceY,
//     targetX,
//     targetY,
//   });

//   console.log("EDGE_ID", id);
//   return (
//     <>
//       <BaseEdge id={Math.random().toString()} path={edgePath} />
//     </>
//   );
// }

// const edgeTypes = {
//   "custom-edge": CustomEdge,
// };

const nodeTypes: NodeTypes = {
  custom: CustomNode,
};

type PorpsType = {
  ast: ASTNode;
};

const ASTtree = ({ ast }: PorpsType) => {
  const [nodes, setNodes, onNodesChange] = useNodesState<any>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<any>([]);

  const onConnect = useCallback(
    (params: Edge | Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  useEffect(() => {
    const newNodes: Node[] = [];
    const newEdges: Edge[] = [];

    const buildGraph = (
      node: ASTNode,
      parentId: string,
      position: { x: number; y: number },
      level: number
    ) => {
      const id = `${node.type}-${node.value}-${Math.random()}`;
      newNodes.push({
        id,
        type: "custom",
        position,
        data: { label: `${node.type}: ${node.value}` },
      });

      if (parentId) {
        newEdges.push({
          id: `${parentId}-${id}`,
          source: parentId,
          sourceHandle: "bottom",
          targetHandle: "top",
          target: id,
        });
      }

      const horizontalSpacing = 200 / (level + 1);
      const verticalSpacing = 100;

      if (node.left) {
        buildGraph(
          node.left,
          id,
          {
            x: position.x - horizontalSpacing,
            y: position.y + verticalSpacing,
          },
          level + 1
        );
      }
      if (node.right) {
        buildGraph(
          node.right,
          id,
          {
            x: position.x + horizontalSpacing,
            y: position.y + verticalSpacing,
          },
          level + 1
        );
      }
    };

    buildGraph(ast, "1", { x: 250, y: 0 }, 0);

    setNodes(newNodes);
    setEdges(newEdges);
  }, [ast]);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        // edgeTypes={edgeTypes}
        fitView
      >
        <Controls />
        <MiniMap />
      </ReactFlow>
    </div>
  );
};

export default ASTtree;
