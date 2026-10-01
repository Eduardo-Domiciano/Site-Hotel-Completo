export const hotel = {
  name: 'Mikhalateia Hotel',
  slogan: 'O único lugar onde você merece estar!',
  phone: '(00) 1234-5678',
  whatsapp: '(00) 9876-5432',
  email: 'comercial@mikhalateia.com.br',
  reservationsEmail: 'reservas@mikhalateia.com.br',
  address: 'Av. ET de Varginha, 01 - Varginha/MG',
  linkedin: 'https://www.linkedin.com/in/eduardo-domiciano/',
  github: 'https://github.com/DuDSTOPIA',
};

export const presentationItems = [
  { title: 'Elegância', image: '/assets/img/Elegancia.jpg' },
  { title: 'Modernidade', image: '/assets/img/Modernidade.jpg' },
  { title: 'Requinte', image: '/assets/img/Requinte.jpg' },
  { title: 'Conforto', image: '/assets/img/Conforto.jpg' },
  { title: 'Comprometimento', image: '/assets/img/Comprometimento.jpg' },
  { title: 'Classe', image: '/assets/img/Classe.jpg' },
  { title: 'Comodidade', image: '/assets/img/Comodidade.jpg' },
  { title: 'Exuberância', image: '/assets/img/Exuberancia.jpg' },
  { title: 'Qualidade', image: '/assets/img/Qualidade.jpg' },
];

export const amenities = [
  { label: 'Conexão wifi', icon: '/assets/img/wifi_icon.svg' },
  { label: 'Restaurante interno', icon: '/assets/img/restaurant.svg' },
  { label: 'Serviço de quarto', icon: '/assets/img/room_service_icon.svg' },
  { label: 'Academia', icon: '/assets/img/gyn_icon.png' },
  { label: 'Piscina', icon: '/assets/img/pool_icon.svg' },
  { label: 'Lavanderia', icon: '/assets/img/laundry_icon.png' },
  { label: 'Garagem', icon: '/assets/img/garage_icon.png' },
  { label: 'Salas para reuniões', icon: '/assets/img/meeting_icon.svg' },
];

export const meetingImages = Array.from(
  { length: 10 },
  (_, i) => `/assets/img/img_slider/sala_reuniao${String(i + 1).padStart(2, '0')}.jpg`,
);

export const loginBackgrounds = [
  '/assets/img/login_slider/hotel_img_bar.jpg',
  '/assets/img/login_slider/hotel_img_breakfest.jpg',
  '/assets/img/login_slider/hotel_img_fiancee.jpg',
  '/assets/img/login_slider/hotel_img_garden.jpg',
  '/assets/img/login_slider/hotel_img_relax.jpg',
  '/assets/img/login_slider/hotel_img_room.jpg',
  '/assets/img/login_slider/hotel_img_rose.jpg',
  '/assets/img/login_slider/hotel_img_towel.jpg',
  '/assets/img/login_slider/hotel_img_wine.jpg',
];

export const reservationPolicies = [
  'Qualquer reserva feita pelo site somente terá validade mediante a apresentação da confirmação enviada por e-mail.',
  'Só é possível selecionar até 03 quartos por reserva. Cada quarto tem um limite de adultos e crianças que não pode ser excedido.',
  'Crianças de até 06 anos são isentas de pagamento da hospedagem.',
  'Reservas sem garantia de no-show serão canceladas automaticamente às 18h na data do check-in.',
  'Cancelamento gratuito até 24 horas antes do check-in. Após esse prazo, cobra-se 50% do valor total.',
  'No check-in é obrigatório documento com foto para adultos e certidão de nascimento para crianças e adolescentes, que só se hospedam acompanhados dos pais, avós ou responsável com documento reconhecido em cartório.',
];

export const careerAreas = [
  'Recepção / Reservas',
  'Administrativo',
  'Governança',
  'Restaurante',
  'Serviços Gerais',
  'TI',
];
