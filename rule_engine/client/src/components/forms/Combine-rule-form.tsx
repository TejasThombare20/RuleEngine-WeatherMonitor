"use client";
import React, { useContext, useEffect, useState } from "react";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { z } from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "../ui/button";
import apiHandler from "@/handlers/api-handler";
import { MultiSelect } from "../MultiSelect";
import { Input } from "../ui/input";
import { useRouter } from "next/navigation";
import { DialogContext } from "../Model";

type Props = {};

const FormSchema = z.object({
  rule_ids: z.array(z.string()).min(2, {
    message: "Please select at least two rule.",
  }),
  name: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  description: z.string().optional(),
});

const CombineRuleForm = (props: Props) => {
  const router = useRouter();
  const { setOpen } = useContext(DialogContext);
  const [rules, setRules] = useState<any>([]);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
      name: "",
      rule_ids: [],
      description: "Combined Rule",
    },
  });

  const options = rules.map((rule: any) => ({
    value: rule?._id,
    label: rule?.name,
  }));

  useEffect(() => {
    const fetchRules = async () => {
      try {
        const rulesData = await apiHandler.get<any>("/rules");
        setRules(rulesData.rules);
        setOpen(false);
      } catch (error) {
        console.error(error);
        setOpen(false);
      }
    };

    fetchRules();
  }, []);

  async function onSubmit(data: z.infer<typeof FormSchema>) {
    console.log("data", data);

    try {
      const combinedRuledata = await apiHandler.post<any>(
        "/rules/combine",
        data
      );
      console.log("combinedRuledata", combinedRuledata);
      router.refresh();
    } catch (error) {}
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Please enter name" {...field} />
              </FormControl>
              <FormDescription>
                Enter New name for combined rule
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Input readOnly disabled placeholder="Description" {...field} />
              </FormControl>
              <FormDescription>Description is auto generated</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="rule_ids"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Rules</FormLabel>
              <Controller
                name="rule_ids"
                control={form.control}
                render={({ field }) => (
                  <MultiSelect
                    options={options}
                    onValueChange={field.onChange}
                    value={field.value}
                    placeholder="Select rules"
                    variant="inverted"
                    animation={2}
                    maxCount={3}
                  />
                )}
              />
              <FormDescription>
                Select the atleast 2 rules for combination
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
};

export default CombineRuleForm;
