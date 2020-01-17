import  PageNotFound from "./components/PageNotFound.vue";
import CarInventory from "./pages/inventory/Car.vue";

export default [
    { path: "/", component: CarInventory },
    { path: "*", component: PageNotFound }
];
