"use client";
import React, { useContext } from "react";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import ApiHandler from "@/handlers/api-handler";
import { DialogContext } from "../Model";
import { Textarea } from "../ui/textarea";
import { useRouter } from "next/navigation";
import { revalidatePath } from "next/cache";

const RuleForm = () => {
  const { setOpen } = useContext(DialogContext);
  const router = useRouter();

  const formSchema = z.object({
    name: z.string().min(2, {
      message: "Username must be at least 2 characters.",
    }),
    description: z.string().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log(values);

    try {
      const ruleData = await ApiHandler.post<any>("/rules", values);
      console.log("ruleData", ruleData);
      setOpen(false);
 
       window.location.reload()

      // setTimeout(() => {
      //   router.refresh();
      // }, 0);
    } catch (error) {
      console.log(error);
      setOpen(false);
    }
  }

  return (
    <div className="w-[500px]">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel> Rule Name </FormLabel>
                <FormControl>
                  <Input placeholder="eg : employee's rule" {...field} />
                </FormControl>
                <FormDescription>Please enter rule title</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Rule String</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="eg : (age > 30 AND department = 'Sales')"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Please enter rule string in specifed format
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit">Submit</Button>
        </form>
      </Form>
    </div>
  );
};

export default RuleForm;
