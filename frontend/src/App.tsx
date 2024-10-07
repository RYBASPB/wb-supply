import {useEffect, useState} from 'react';
import './App.css';
import {
    CreateApiClient,
    DeleteApiClient,
    DeleteApiKey,
    GetAllApiKeys,
    GetWildberriesData,
    SaveApiKey
} from "../wailsjs/go/main/App";
import {sqlite} from "../wailsjs/go/models";
import {Button, ButtonGroup, Input, Select, SelectItem} from "@nextui-org/react";

type NewApiKey = Partial<sqlite.ApiKey>;

function App() {
    type Stats = {
        [article: string]:
            {
                [warehouse: string]:
                    {
                        [group: string]: number
                    }
            }
    }
    const [newKey, setNewKey] = useState<NewApiKey>({
        name: "",
        api_key: ""
    });
    const [fromWb, setFromWb] = useState<{articles: {}, stats: Stats}>({
        articles: {},
        stats: {}
    });
    const [apiKeysList, setApiKeysList] = useState<sqlite.ApiKey[]>([]);
    const [activeKey, setActiveKey] = useState<number>(0);
    const [startDate, setStartDate] = useState<string>("");
    const [isLoading, setIsLoading] = useState(false);

    const updateDate = (value: any) => {
        const parsedDate = new Date(value);
        if (parsedDate) {
            const formatedDate = parsedDate.toISOString();
            setStartDate(formatedDate);
        }
    }
    const updateName = (e: any) => setNewKey(prevState => ({...prevState, name: e.target.value}));
    const addApiKey = (e: any) => setNewKey(prevState => ({...prevState, api_key: e.target.value}));
    const deleteKey = (e: any) => deleteKeyFromDb(e.target.getAttribute("data-id"))

    function addApiKeyToDb() {
        if (newKey.api_key && newKey.name) {
            SaveApiKey(newKey.api_key, newKey.name).then(() => loadKeys());
        }
    }

    useEffect(() => {
        loadKeys().catch(console.error);
    }, [])

    async function loadKeys() {
        const keys = await GetAllApiKeys();
        setApiKeysList(keys);
    }

    function deleteKeyFromDb(id: string) {
        DeleteApiKey(id).catch(e => console.error(e));
        loadKeys().catch(console.error);
    }

    async function getData() {
        if (activeKey && startDate.length > 0) {
            setIsLoading(true);
            let {articles, stats} = await GetWildberriesData(startDate);
            articles = JSON.parse(articles);
            stats = JSON.parse(stats);
            // @ts-ignore
            setFromWb({articles, stats});
            setIsLoading(false);
        }
    }

    async function activateKey(e: any) {
        const keyIdToFind = parseInt(e.target.value);
        const key = apiKeysList.find(key => key.id === keyIdToFind);

        if (!key) {
            console.error("Ключ не найден")
            return
        }

        if (activeKey == keyIdToFind) {
            //TODO add error if no key
            console.error("Ключ уже активирован")
            return
        }
        if (activeKey != keyIdToFind) {
            await deleteClient();
        }
        await CreateApiClient(key.api_key);
        setActiveKey(key.id);
    }
    
    async function deleteClient() {
        await DeleteApiClient();
        setActiveKey(0);
    }

    return (
        <div id="App" className="flex flex-col gap-4 items-center m-4">
            <h1 className="text-xl">Расчет поставок Wildberries</h1>
            <div className="flex items-center gap-2">
                <Input
                    isRequired
                    isClearable
                    label="Дата начала"
                    type="date"
                    onValueChange={updateDate}
                />
                <ButtonGroup>
                    <Button onPress={loadKeys}>Обновить ключи</Button>
                    <Button onPress={getData}>Запрос данных</Button>
                </ButtonGroup>
            </div>
            <div>
                {/*TODO table presentation*/}
                {isLoading
                    ? "Loading"
                    : Object.keys(fromWb.stats).map(key => {
                        let result: string[] = [];
                        const warehouses = fromWb.stats[key]
                        const warehouseNames = Object.keys(warehouses);
                        warehouseNames.map(warehouse => {
                            const {sales, stocks, supply} = warehouses[warehouse]
                            result.push(`${key}: ${warehouse}: sales: ${sales}, stocks: ${stocks}, supply: ${supply}`)
                        })
                        return <div>
                            {result.map(str => <div>{str}</div>)}
                        </div>
                    })
                }
            </div>
            <Select
                items={apiKeysList}
                selectedKeys={[activeKey.toString()]}
                onChange={activateKey}
            >
                {key => <SelectItem key={key.id}>
                    {key.name}
                </SelectItem>}
            </Select>
            <Button size="sm" color="danger" onPress={deleteKey} data-id={activeKey}>Удалить</Button>

            <div className="flex flex-col max-w-lg gap-4">
                    <Input id="name" label="Название ключа" className="input" onChange={updateName} autoComplete="off" name="input"
                           type="text"/>
                    <Input id="apiKey" label="Ключ API" className="input" onChange={addApiKey} autoComplete="off" name="apiKey"
                           type="text"/>
                <Button className="btn" onPress={addApiKeyToDb}>Добавить ключ</Button>
            </div>
        </div>
    )
}

export default App
