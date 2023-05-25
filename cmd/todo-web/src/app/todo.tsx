import { useState, useEffect } from "preact/hooks"

type TodoItem = {
    id: number
    description: string
    completed: boolean
}

export function Todo({ ipn }: { ipn: IPN; }) {
    const [todoList, setTodoList] = useState([] as TodoItem[]);
    const [newTodoDescription, setNewTodoDescription] = useState('');

    useEffect(() => {
        ipn.fetch('http://todo-server')
            .then((response) => response.text())
            .then((data) => {
                // parse this JSON and return the list of items
                let items: TodoItem[] = [];
                JSON.parse(data).map((item: any) => {
                    items.push({
                        id: item.id,
                        description: item.description,
                        completed: item.completed ? true : false,   // convert to boolean
                    });
                });

                setTodoList(items);
            })
            .catch((error) => console.error(error));
    }, []);

    // send a POST request to the server to add a new todo item
    const handleAddTodo = () => {
        const newTodoItem = { description: newTodoDescription, completed: false };
        console.log("request", JSON.stringify(newTodoItem));
        ipn.post('http://todo-server', 'application/io.goswagger.examples.todo-list.v1+json', JSON.stringify(newTodoItem))
            .then((response) => {
                console.log(response);
                return response.text();
            })
            .then((data) => {
                // add the new todo item to the todo list
                const newTodoItem = JSON.parse(data);
                setTodoList([...todoList, newTodoItem]);
                setNewTodoDescription('');
            })
            .catch((error) => console.error(error));
    };

    // send a POST request to the server when a checkbox is clicked
    const handleCheckboxChange = (item: TodoItem) => {
        const updatedItem = { ...item, completed: !item.completed };
        // remove id field from the updated item
        const updateItemBody = { description: updatedItem.description, completed: updatedItem.completed };

        const url = `http://todo-server/${updatedItem.id}`;
        ipn.put(url, "application/io.goswagger.examples.todo-list.v1+json", JSON.stringify(updateItemBody))
            .then((response) => {
                console.log(response);
                return response.text()
            })
            .then((data) => {
                // update the todo list with the updated item
                const updatedList = todoList.map((item) => {
                    if (item.id === updatedItem.id) {
                        return updatedItem;
                    } else {
                        return item;
                    }
                });
                setTodoList(updatedList);
            })
            .catch((error) => console.error(error));
    };

    // print each item in the list
    return (
        <div className="container mx-auto px-4">
            <div className="flex flex-col items-center justify-items-center">
                {todoList.map((item) => {
                    return (
                        <div key={item.id} className="flex flex-row items-center justify-items-center my-2">
                            <input type="checkbox" className="mr-2" checked={item.completed} onChange={() => handleCheckboxChange(item)} />
                            <div className={`text-left ${item.completed ? 'line-through' : ''}`}>{item.description}</div>
                        </div>
                    );
                })}
                <div className="flex flex-row items-center justify-items-center my-2">
                    <input type="text" className="border border-gray-400 rounded py-2 px-4 mr-2" value={newTodoDescription} onChange={(event) => {
                        if (event.target instanceof HTMLInputElement) {
                            setNewTodoDescription(event.target.value)
                        }
                    }} />
                    <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" onClick={handleAddTodo}>Add Todo</button>
                </div>
            </div>
        </div>
    );
}