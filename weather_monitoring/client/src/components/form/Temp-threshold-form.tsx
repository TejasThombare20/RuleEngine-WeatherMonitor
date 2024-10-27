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
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import apiHandler from "@/handlers/apiHandler";
import { useRouter } from "next/navigation";
import { DialogContext } from "../Temp-threshold";
import { useToast } from "@/hooks/use-toast";

type Props = {};

const TempThresholdForm = (props: Props) => {
  const router = useRouter();
  const { toast } = useToast();
  const { setOpen } = useContext(DialogContext);

  const formSchema = z.object({
    email: z.string().email(),
    consecutiveAlerts: z.coerce.number().gte(2, {
      message: "Min value at least 2",
    }),
    Mumbai: z.coerce.number(),
    Delhi: z.coerce.number(),
    Bengaluru: z.coerce.number(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    const RequestBody = {
      email: values.email,
      consecutiveAlerts: values.consecutiveAlerts,
      Thres_temperatue: {
        Mumbai: values.Mumbai,
        Delhi: values.Delhi,
        Bengaluru: values.Bengaluru,
      },
    };

    try {
      await apiHandler.post("/createuser", RequestBody);
      setOpen(false);
      toast({ title: "User created successfully" });
    } catch (error) {
      console.log(error);
      setOpen(false);
      toast({
        title: "failed to create a user,Please try with another email ",
        variant: "destructive",
      });
    }
  }

  return (
    <div className="w-full">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input placeholder="Email" {...field} />
                </FormControl>
                <FormDescription>
                  You will receive notification this email
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="consecutiveAlerts"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Consecutive occurrences</FormLabel>
                <FormControl>
                  <Input
                    type="number"
                    placeholder="Enter the number of consecutive occurrences"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  If the temperature exceeds the threshold for the above
                  specified number of consecutive times, you will receive a
                  notification
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="font-semibold">
            Enter the temperatue threashold value for following clities
          </div>
          <div className="w-full flex justify-between items-center ">
            <FormField
              control={form.control}
              name="Mumbai"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>For Mumbai</FormLabel>
                  <FormControl>
                    <Input
                      {...field}
                      type="number"
                      placeholder="in Celcius"
                      {...field}
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="Delhi"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>For Delhi</FormLabel>
                  <FormControl>
                    <Input type="number" placeholder="in Celcius" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <FormField
            control={form.control}
            name="Bengaluru"
            render={({ field }) => (
              <FormItem>
                <FormLabel>For Bengluru</FormLabel>
                <FormControl>
                  <Input type="number" placeholder="in Celcius" {...field} />
                </FormControl>
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

export default TempThresholdForm;
