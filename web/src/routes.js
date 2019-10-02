import Releases from './components/Releases.vue'
import NotFoundComponent from "./components/NotFoundComponent";
import HomeComponent from "./components/HomeComponent";
import ReleaseComponent from "./components/ReleaseComponent";

const routes = [
  { path: '/', component: HomeComponent },
  { path: '/releases', component: Releases },
  { path: '/releases/:name', component: ReleaseComponent},
  { path: '*', component: NotFoundComponent }
];

export default routes
