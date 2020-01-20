import  PageNotFound from "./pages/common/PageNotFound.vue";
import CarInventory from "./pages/inventory/Car.vue";

export default [
    { path: "/", redirect: "/inventory/cars" },
    { path: "/inventory/cars", component: CarInventory },
    { path: "*", component: PageNotFound }
];
