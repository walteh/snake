import { useEffect, useMemo } from "react";
import { Fragment, useState } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { swails } from "../wailsjs/go/models";
import { CopyBlock } from "react-code-blocks";
import clsx from "clsx";
import ReactJson from "@microlink/react-json-view";

export default function Modal({
	result,
}: {
	result: swails.WailsHTMLResponse | undefined;
}) {
	const [open, setOpen] = useState(true);

	useEffect(() => {
		if (result) {
			setOpen(true);
		}
	}, [result]);

	return (
		<Transition.Root show={open} as={Fragment}>
			<Dialog
				as="div"
				className={clsx("relative", "z-10")}
				onClose={setOpen}
			>
				<Transition.Child
					as={Fragment}
					enter="ease-out duration-300"
					enterFrom="opacity-0"
					enterTo="opacity-100"
					leave="ease-in duration-200"
					leaveFrom="opacity-100"
					leaveTo="opacity-0"
				>
					<div
						className={clsx(
							"fixed",
							"inset-0",
							"bg-gray-500",
							"bg-opacity-75",
							"transition-opacity"
						)}
					/>
				</Transition.Child>

				<div
					className={clsx(
						"fixed",
						"inset-0",
						"z-10",
						"w-screen",
						"overflow-y-auto"
					)}
				>
					<div
						className={clsx(
							"flex",
							// "h-full",
							"items-start",
							"justify-center",
							// "justify-center",
							"p-10",
							"text-center",
							// "sm:items-center",
							// "h-screen",
							"sm:p-0",
							"m-10"
						)}
					>
						<Transition.Child
							as={Fragment}
							enter="ease-out duration-300"
							enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
							enterTo="opacity-100 translate-y-0 sm:scale-100"
							leave="ease-in duration-200"
							leaveFrom="opacity-100 translate-y-0 sm:scale-100"
							leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
						>
							<Dialog.Panel
								className={clsx(
									"w-10/12",
									// "h-full",
									"transform",
									"rounded-lg",
									"bg-white",
									"px-4",
									"pb-4",
									"pt-5",
									"text-left",
									"shadow-xl",
									"transition-all"
								)}
							>
								<div>
									<div className={clsx()}>
										{result && <Tabs arg={result} />}
									</div>
								</div>
							</Dialog.Panel>
						</Transition.Child>
					</div>
				</div>
			</Dialog>
		</Transition.Root>
	);
}

function Table({ data }: { data: string[][] }) {
	return (
		<div className="px-4 sm:px-6 lg:px-8">
			<div className="mt-8 flow-root">
				<div className="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
					<div className="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
						<div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
							<table className="min-w-full divide-y divide-gray-300">
								<thead className="bg-gray-50">
									<tr>
										{data[0].map((header) => (
											<th
												key={header}
												scope="col"
												className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900"
											>
												{header}
											</th>
										))}
									</tr>
								</thead>
								<tbody className="divide-y divide-gray-200 bg-white">
									{data.map(
										(person, i) =>
											i !== 0 && (
												<tr key={i}>
													{person.map((vals, i) => (
														<td
															key={i}
															className="px-3 py-4 text-sm text-gray-500"
														>
															{vals}
														</td>
													))}
												</tr>
											)
									)}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

const Tabs = ({ arg }: { arg: swails.WailsHTMLResponse }) => {
	const [current, setCurrent] = useState(0);

	const tabs = useMemo(() => {
		let tabs: { name: string; content: JSX.Element }[] = [];

		if (arg.json) {
			tabs.push({
				name: "JSON",
				content: (
					<ReactJson
						src={arg.json}
						theme={"google"}
						style={{
							padding: "1rem",
							borderRadius: "0.5rem",
						}}
					/>
				),
			});
		}

		if (arg.table) {
			tabs.push({ name: "Table", content: <Table data={arg.table} /> });
		}

		if (arg.text) {
			tabs.push({
				name: "Text",
				content: (
					<code className="whitespace-pre-wrap">{arg.text}</code>
				),
			});
		}

		return tabs;
	}, [arg]);

	return (
		<div>
			<div className="sm:hidden">
				<label htmlFor="tabs" className="sr-only">
					Select a tab
				</label>
				{/* Use an "onChange" listener to redirect the user to the selected tab URL. */}
				<select
					id="tabs"
					name="tabs"
					className="block w-full rounded-md border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
					defaultValue={tabs[0].name}
				>
					{tabs.map((tab) => (
						<option key={tab.name}>{tab.name}</option>
					))}
				</select>
			</div>
			<div className="hidden sm:block">
				<nav className="flex space-x-4" aria-label="Tabs">
					{tabs.map((tab, i) => (
						<a
							key={tab.name}
							onClick={() => {
								setCurrent(i);
							}}
							className={clsx(
								tabs[current].name === tab.name
									? "bg-indigo-100 text-indigo-700"
									: "text-gray-500 hover:text-gray-700",
								"rounded-md px-3 py-2 text-sm font-medium"
							)}
							aria-current={tabs[current] ? "page" : undefined}
						>
							{tab.name}
						</a>
					))}
				</nav>
			</div>
			<div
				className={clsx(
					"mt-3",
					"border-t",
					"border-gray-200",
					"pt-4",
					"sm:pt-6",
					"overflow-y-auto"
				)}
			>
				{tabs[current].content}
			</div>
		</div>
	);
};
