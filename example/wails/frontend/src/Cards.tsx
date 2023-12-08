import { useCallback, useEffect, useMemo, useState } from "react";
import { Menu, Transition } from "@headlessui/react";
import { EllipsisHorizontalIcon } from "@heroicons/react/20/solid";
import { swails } from "../wailsjs/go/models";
import { InputsFor, Run, RunWithWriter } from "../wailsjs/go/swails/WailsSnake";
import Input from "./Input";
import Modal from "./Modal";

export function Card({ command }: { command: swails.WailsCommand }) {
	const [allInputs, setAllInputs] = useState<swails.WailsInput[]>();

	useEffect(() => {
		InputsFor(command).then((result) => {
			const obj: Map<string, swails.WailsInput> = new Map();
			console.log("Inputs: ", { result });
			for (let i = 0; i < result.length; i++) {
				obj.set(result[i].name, result[i]);
			}
			setAllInputs(result);
		});
	}, [command]);

	const [writer, setWriter] = useState<swails.WailsWriter>();

	const caller = useCallback(function execute(cmd: swails.WailsCommand) {
		const random = Math.random().toString(36).substring(7);
		const wrt = new swails.WailsWriter();
		wrt.name = random;
		setWriter(wrt);
		console.log("writer: ", wrt);
		RunWithWriter(random, cmd).then((result) => {
			// console.log(result);
			// setWriter(result);

			console.log("result: ", result);
		});
	}, []);

	return (
		<li
			key={command.name}
			className="overflow-hidden rounded-xl border border-gray-200"
		>
			<div className="flex flex-col justify-start items-start gap-x-4 border-b border-gray-900/5 bg-gray-50 p-6">
				<div className="flex justify-between items-center w-full">
					<div className="text-md font-semibold leading-6 text-gray-900">
						{command.name}
					</div>
					<button
						type="button"
						onClick={() => caller(command)}
						className="rounded-md bg-indigo-50 px-3 py-2 text-sm font-semibold text-indigo-600 shadow-sm hover:bg-indigo-100"
					>
						execute
					</button>
				</div>

				<p className="mt-2 text-sm text-gray-500">
					{command.description}
				</p>
			</div>

			<dl className="-my-3 divide-y divide-gray-100 px-6 py-4 text-sm leading-6">
				{allInputs?.map((input) => (
					<div
						className="flex justify-between gap-x-2 py-3"
						key={input.name}
					>
						<Input arg={input} />
					</div>
				))}
			</dl>

			<Modal result={writer} />
		</li>
	);
}

export default function Cards({
	commands,
}: {
	commands: swails.WailsCommand[];
}) {
	return (
		<ul
			role="list"
			className="grid grid-cols-1 gap-x-6 gap-y-8 lg:grid-cols-3 xl:gap-x-8 p-20"
		>
			{commands.map((client) => (
				<Card key={client.name} command={client} />
			))}
		</ul>
	);
}
