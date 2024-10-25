"use client";
import React, { useState } from "react";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Button } from "../ui/button";
import { string, z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Textarea } from "../ui/textarea";
import apiHandler from "@/handlers/api-handler";

type Props = { rule_id: string };

const formSchema = z.object({
  jsonInput: z.string().refine((value) => {
    try {
      JSON.parse(value);
      return true;
    } catch {
      return false;
    }
  }, "Invalid JSON format"),
});

const Verifyruleform = ({ rule_id }: Props) => {
  const [validationResult, setValidationResult] = useState<boolean | null>(
    null
  );

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      jsonInput: "{ \n}",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    try {
      console.log("values", values);
      const parsedJson = JSON.parse(values.jsonInput);
      console.log("parsedJson", parsedJson);

      const requestData = {
        rule_id,
        data: parsedJson,
      };

      const RuleVerifyData = await apiHandler.post<any>(
        "/evaluate",
        requestData
      );

      const isPassed = RuleVerifyData?.result?.result;

      setValidationResult(isPassed);
      //   toast({
      //     title: "Validation Successful",
      //     description: "The input is valid JSON.",
      //   });
    } catch (error) {
      setValidationResult(false);
      //   toast({
      //     title: "Validation Failed",
      //     description: "The input is not valid JSON.",
      //     variant: "destructive",
      //   });
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="jsonInput"
          render={({ field }) => (
            <FormItem>
              <FormLabel>JOSN input</FormLabel>
              <FormControl>
                <Textarea
                  placeholder='(e.g. {"name": "John", "age": 30})'
                  className="min-h-[100px] "
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Enter a valid JSON data for rule evaluation.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Validate JSON</Button>
        {validationResult && (
          <div
            className={
              validationResult === true ? "text-green-600" : "text-red-600"
            }
          >
            {validationResult === true
              ? "Input example is valid for this rule"
              : "Input example is not valid for this rule"}
          </div>
        )}
      </form>
    </Form>
  );
};

export default Verifyruleform;
