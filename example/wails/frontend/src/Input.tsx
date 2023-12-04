import { useCallback, useEffect, useMemo, useState } from "react";
import { UpdateInput, OptionsForEnum } from "../wailsjs/go/swails/WailsSnake";
import { swails } from "../wailsjs/go/models";

import { Switch } from "@headlessui/react";
import clsx from "clsx";

const BoolInput = ({ arg }: { arg: swails.WailsInput }) => {
	const [value, setValue] = useState<swails.WailsInput>(arg);

	const enabled = useMemo(() => {
		return value.value;
	}, [value]);

	const setEnabled = useCallback(
		(next: boolean) => {
			const newInput = value;
			newInput.value = next;
			UpdateInput(newInput).then((result) => {
				console.log({ result });
				setValue(result);
			});
		},
		[value.value]
	);

	return (
		<Switch.Group as="div" className="flex items-center">
			<Switch
				checked={enabled}
				onChange={setEnabled}
				className={clsx(
					enabled ? "bg-indigo-600" : "bg-gray-200",
					"relative",
					"inline-flex",
					"h-6",
					"w-11",
					"flex-shrink-0",
					"cursor-pointer",
					"rounded-full",
					"border-2",
					"border-transparent",
					"transition-colors",
					"duration-200",
					"ease-in-out",
					"focus:outline-none",
					"focus:ring-2",
					"focus:ring-indigo-600",
					"focus:ring-offset-2"
				)}
			>
				<span className="sr-only">Use setting</span>
				<span
					className={clsx(
						enabled ? "translate-x-5" : "translate-x-0",
						"pointer-events-none",
						"relative",
						"inline-block",
						"h-5",
						"w-5",
						"transform",
						"rounded-full",
						"bg-white",
						"shadow",
						"ring-0",
						"transition",
						"duration-200",
						"ease-in-out"
					)}
				>
					<span
						className={clsx(
							enabled
								? "opacity-0 duration-100 ease-out"
								: "opacity-100 duration-200 ease-in",
							"absolute",
							"inset-0",
							"flex",
							"h-full",
							"w-full",
							"items-center",
							"justify-center",
							"transition-opacity"
						)}
						aria-hidden="true"
					>
						<svg
							className="h-3 w-3 text-gray-400"
							fill="none"
							viewBox="0 0 12 12"
						>
							<path
								d="M4 8l2-2m0 0l2-2M6 6L4 4m2 2l2 2"
								stroke="currentColor"
								strokeWidth={2}
								strokeLinecap="round"
								strokeLinejoin="round"
							/>
						</svg>
					</span>
					<span
						className={clsx(
							enabled
								? "opacity-100 duration-200 ease-in"
								: "opacity-0 duration-100 ease-out",
							"absolute",
							"inset-0",
							"flex",
							"h-full",
							"w-full",
							"items-center",
							"justify-center",
							"transition-opacity"
						)}
						aria-hidden="true"
					>
						<svg
							className="h-3 w-3 text-indigo-600"
							fill="currentColor"
							viewBox="0 0 12 12"
						>
							<path d="M3.707 5.293a1 1 0 00-1.414 1.414l1.414-1.414zM5 8l-.707.707a1 1 0 001.414 0L5 8zm4.707-3.293a1 1 0 00-1.414-1.414l1.414 1.414zm-7.414 2l2 2 1.414-1.414-2-2-1.414 1.414zm3.414 2l4-4-1.414-1.414-4 4 1.414 1.414z" />
						</svg>
					</span>
				</span>
			</Switch>
			<Switch.Label as="span" className="ml-3 text-sm">
				<span className=" text-gray-500  font-bold">{value.name}</span>{" "}
			</Switch.Label>
		</Switch.Group>
	);
};

const StringInput = ({ arg }: { arg: swails.WailsInput }) => {
	const [value, setValue] = useState<swails.WailsInput>(arg);

	return (
		<div
			className={clsx(
				"relative",
				"rounded-md",
				"px-3",
				"pb-1.5",
				"pt-2.5",
				"shadow-sm",
				"ring-1",
				"ring-inset",
				"ring-gray-300",
				"focus-within:ring-2",
				"focus-within:ring-indigo-600"
			)}
		>
			<label
				htmlFor={value.name}
				className="absolute -top-2.5 left-2 inline-block bg-white px-1 text-sm font-medium"
			>
				<span className="text-gray-500  font-bold">{arg.name}</span>
			</label>
			<input
				id="name"
				type="text"
				name={value.name}
				value={value.value}
				onChange={(next) => {
					const newInput = value;
					newInput.value = next.target.value;
					UpdateInput(newInput).then((result) => {
						console.log({ result });
						setValue(result);
					});
				}}
				className="block w-full border-0 p-0 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
				placeholder={value.name}
			/>
		</div>
	);
};

const EnumInput = ({ arg }: { arg: swails.WailsInput }) => {
	const [value, setValue] = useState<swails.WailsInput>(arg);

	const [options, setOptions] = useState<string[]>([]);

	useEffect(() => {
		OptionsForEnum(arg).then((result) => {
			setOptions(result);
		});
	}, [arg]);

	return (
		<div className="mt-2 flex rounded-md shadow-sm relative ">
			<label
				htmlFor={value.name}
				className="absolute -top-1 left-2 inline-block bg-white px-1 text-sm font-medium"
			>
				<span className="text-gray-500  font-bold">{arg.name}</span>
			</label>
			<select
				id={value.name}
				name={value.name}
				className="mt-2 block rounded-md border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6 w-full"
				defaultValue="Canada"
				onChange={(next) => {
					const newInput = value;
					newInput.value = next.target.value;
					UpdateInput(newInput).then((result) => {
						console.log({ result });
						setValue(result);
					});
				}}
			>
				{options?.map((option) => {
					return (
						<option key={option} value={option}>
							{option}
						</option>
					);
				})}
			</select>
		</div>
	);
};

const Input = ({ arg }: { arg: swails.WailsInput | undefined }) => {
	switch (arg?.type) {
		case "string":
			return <StringInput arg={arg} />;
		case "bool":
			return <BoolInput arg={arg} />;
		case "enum":
			return <EnumInput arg={arg} />;
		default:
			return <div>unknown input type</div>;
	}
};

export default Input;
